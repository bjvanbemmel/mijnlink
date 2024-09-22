package service

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/bjvanbemmel/mijnlink/utils"
)

var (
	ErrCorruptData = errors.New("could not read data, possibly corrupt")
	ErrNotFound    = errors.New("resource not found")
)

type URLService struct {
	KeyLimit int
	File     *os.File
	Mutex    *sync.Mutex
}

func (s URLService) SaveUrl(url string) (string, error) {
	if key, _ := s.GetKeyByURL(url); key != "" {
		return key, nil
	}

	var key string
	for {
		key = utils.Key(s.KeyLimit)
		found, err := s.GetKeyByKey(key)
		if err != nil && !errors.Is(err, ErrNotFound) {
			return "", err
		}

		if found == "" {
			break
		}
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	line := fmt.Sprintf("%s=%s\n", key, url)
	_, err := s.File.WriteString(line)

	return key, err
}

func (s URLService) GetURL(callback func(key string, url string) bool) (string, string, error) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.File.Seek(0, 0)

	scanner := bufio.NewScanner(s.File)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		split := strings.SplitN(line, "=", 2)
		if len(split) < 2 {
			return "", "", ErrCorruptData
		}

		key := split[0]
		url := split[1]

		if !callback(key, url) {
			continue
		}

		return key, url, nil
	}

	return "", "", ErrNotFound
}

func (s URLService) GetURLByKey(key string) (string, error) {
	_, url, err := s.GetURL(func(k string, _ string) bool {
		return k == key
	})

	return url, err
}

func (s URLService) GetKeyByKey(key string) (string, error) {
	key, _, err := s.GetURL(func(k string, _ string) bool {
		return k == key
	})

	return key, err
}

func (s URLService) GetKeyByURL(url string) (string, error) {
	key, _, err := s.GetURL(func(_ string, u string) bool {
		return u == url
	})

	return key, err
}
