package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var store map[string]string = make(map[string]string)

func main() {
  r := chi.NewRouter()
  r.Use(middleware.Recoverer)
  r.Use(middleware.Logger)

  r.Get("/{key}", func(w http.ResponseWriter, r *http.Request) {
    key := chi.URLParam(r ,"key")
    slog.Info("", "key", key)
    url, err := GetUrl("test")
    if err != nil {
      slog.Error("err", err)
      return
    }
    fmt.Println(store)

    http.Redirect(w, r, fmt.Sprintf("https://%s", url), http.StatusMovedPermanently)
  })

  r.Post("/{url}", func(w http.ResponseWriter, r *http.Request) {
    url := chi.URLParam(r, "url")

    key, err := NewUrl(url)
    if err != nil {
      slog.Error("err", err)
      return
    }

    w.Write([]byte(key))
  })

  http.ListenAndServe(":80", r)
}

func NewUrl(url string) (string, error) {
  store["test"] = url
  fmt.Println(store)

  return "test", nil
}

func GetUrl(key string) (string, error) {
  return store[key], nil
}
