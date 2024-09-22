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

type IndexService struct {
	KeyLimit int
	File     *os.File
	Mutex    *sync.Mutex
}

var (
	ErrCorruptData = errors.New("could not read data, possibly corrupt")
	ErrNotFound    = errors.New("resource not found")
)

func (s IndexService) SaveValue(value string) (string, error) {
	if key, _ := s.GetKeyByValue(value); key != "" {
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
	line := fmt.Sprintf("%s=%s\n", key, value)
	_, err := s.File.WriteString(line)

	return key, err
}

func (s IndexService) GetValue(callback func(key string, value string) bool) (string, string, error) {
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
		value := split[1]

		if !callback(key, value) {
			continue
		}

		return key, value, nil
	}

	return "", "", ErrNotFound
}

func (s IndexService) GetValueByKey(key string) (string, error) {
	_, value, err := s.GetValue(func(k string, _ string) bool {
		return k == key
	})

	return value, err
}

func (s IndexService) GetKeyByKey(key string) (string, error) {
	key, _, err := s.GetValue(func(k string, _ string) bool {
		return k == key
	})

	return key, err
}

func (s IndexService) GetKeyByValue(value string) (string, error) {
	key, _, err := s.GetValue(func(_ string, v string) bool {
		return v == value
	})

	return key, err
}
