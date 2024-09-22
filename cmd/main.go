package main

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/bjvanbemmel/mijnlink/controller"
	"github.com/bjvanbemmel/mijnlink/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

const (
	TIMEOUT_DURATION = time.Second * 3
	URL_FILE_PATH    = ".index"
)

var (
	ErrPrefixNotSet = errors.New("url prefix has not been set")
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Timeout(TIMEOUT_DURATION))
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST"},
	}))

	environment := os.Getenv("ENVIRONMENT")
	if environment == "debug" {
		r.Use(middleware.Logger)
	}

	kl := os.Getenv("KEY_LENGTH")
	keyLength, err := strconv.Atoi(kl)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(URL_FILE_PATH, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}

	urlsrv := service.URLService{
		KeyLimit: keyLength,
		File:     file,
		Mutex:    &sync.Mutex{},
	}

	prefix := os.Getenv("URL_PREFIX")
	if prefix == "" {
		panic(ErrPrefixNotSet)
	}

	urlctrl := controller.URLController{
		URLService: urlsrv,
		URLPrefix:  prefix,
	}
	urlctrl.InitRoutes(r)

	http.ListenAndServe(":80", r)
}
