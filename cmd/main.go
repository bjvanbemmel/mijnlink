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
	ErrURLPrefixNotSet  = errors.New("url prefix has not been set")
	ErrFilePrefixNotSet = errors.New("file prefix has not been set")
	ErrInvalidSizeLimit = errors.New("invalid upload size limit")
	IndexMutex          = &sync.Mutex{}
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

	keyLength, err := strconv.Atoi(os.Getenv("KEY_LENGTH"))
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(URL_FILE_PATH, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		panic(err)
	}

	indexsvc := service.IndexService{
		KeyLimit: keyLength,
		File:     file,
		Mutex:    &sync.Mutex{},
	}

	urlsvc := service.URLService{
		IndexService: indexsvc,
	}

	urlPrefix := os.Getenv("URL_PREFIX")
	if urlPrefix == "" {
		panic(ErrURLPrefixNotSet)
	}

	urlctrl := controller.URLController{
		URLService: urlsvc,
		URLPrefix:  urlPrefix,
	}
	urlctrl.InitRoutes(r)

	filesvc := service.FileService{
		IndexService: indexsvc,
	}

	limit, err := strconv.Atoi(os.Getenv("UPLOAD_SIZE_LIMIT"))
	if err != nil {
		panic(ErrInvalidSizeLimit)
	}

	filePrefix := os.Getenv("FILE_PREFIX")
	if filePrefix == "" {
		panic(ErrFilePrefixNotSet)
	}

	filectrl := controller.FileController{
		URLPrefix:       filePrefix,
		FileService:     filesvc,
		UploadSizeLimit: limit,
	}
	filectrl.InitRoutes(r)

	http.ListenAndServe(":80", r)
}
