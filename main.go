package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	TIMEOUT_DURATION = time.Second * 3
	STORAGE_FILENAME = ".data"
)

var (
	keyLength            int
	errInvalidRequest    = errors.New("invalid request body given")
	errCouldNotAddToFile = errors.New("something went wrong while saving this url")
	errKeyNotFound       = errors.New("the given key does not exist")
	errUrlNotFound       = errors.New("the given url does not exist")
)

type Response struct {
	Value string `json:"value"`
}

type Request struct {
	URL string `json:"url"`
}

func main() {
	k := os.Getenv("KEY_LENGTH")

	var err error
	keyLength, err = strconv.Atoi(k)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Timeout(TIMEOUT_DURATION))
	r.Use(middleware.Recoverer)

	r.Post("/", NewUrlHandler)
	r.Get("/{key}", GetUrlHandler)

	http.ListenAndServe(":80", r)
}

func NewUrlHandler(w http.ResponseWriter, r *http.Request) {
	var request Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		Error(w, errInvalidRequest, http.StatusBadRequest)
		return
	}

	key, err := AddToFile(request.URL)
	if err != nil {
		Error(w, errCouldNotAddToFile, http.StatusInternalServerError)
		return
	}

	OK(w, key)
}

func GetUrlHandler(w http.ResponseWriter, r *http.Request) {
	_, url, err := FindInFile(chi.URLParam(r, "key"))
	if err != nil {
		Error(w, err, http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusPermanentRedirect)
	return
}

func OK(w http.ResponseWriter, msg string) {
	res := Response{
		Value: msg,
	}

	raw, _ := json.Marshal(res)
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", raw)
}

func Error(w http.ResponseWriter, err error, status int) {
	res := Response{
		Value: err.Error(),
	}

	raw, _ := json.Marshal(res)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, "%s", raw)
}

type Storage struct {
	file  *os.File
	mutex *sync.Mutex
}

func AddToFile(url string) (string, error) {
	file, err := os.OpenFile(STORAGE_FILENAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return "", err
	}
	defer file.Close()

	key, _, _ := FindInFileByUrl(url)
	if key != "" {
		return key, nil
	}

	// Generate new random string in case of duplicate
	unique := false
	key = ""
	for unique == false {
		key = RandomString()
		_, url, _ := FindInFile(key)
		unique = url == ""
	}

	_, err = file.WriteString(fmt.Sprintf("%s=%s\n", key, url))
	return key, err
}

func FindInFile(key string) (string, string, error) {
	file, err := os.Open(STORAGE_FILENAME)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		substrings := strings.SplitN(scanner.Text(), "=", 2)
		if key == substrings[0] {
			return substrings[0], substrings[1], nil
		}
	}

	return "", "", errKeyNotFound
}

func FindInFileByUrl(url string) (string, string, error) {
	file, err := os.Open(STORAGE_FILENAME)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		substrings := strings.SplitN(scanner.Text(), "=", 2)
		if url == substrings[1] {
			return substrings[0], substrings[1], nil
		}
	}

	return "", "", errUrlNotFound
}

func RandomString() string {
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890")

	s := make([]rune, keyLength)
	for i := range s {
		s[i] = chars[rand.Intn(len(chars))]
	}

	return string(s)
}
