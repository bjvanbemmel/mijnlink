package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/bjvanbemmel/mijnlink/response"
	"github.com/bjvanbemmel/mijnlink/service"
	"github.com/go-chi/chi/v5"
)

var (
	ErrInvalidKey     = errors.New("given key is invalid")
	ErrInvalidRequest = errors.New("invalid request body")
	ErrInvalidUrl     = errors.New("invalid url given")
	ErrUrlNotAllowed  = errors.New("the given url is not allowed")
	ErrObfuscation    = errors.New("something went wrong")
)

type URLController struct {
	URLPrefix  string
	URLService service.URLService
}

func (c URLController) InitRoutes(r *chi.Mux) {
	r.Post("/url", c.saveURL)
	r.Get("/url/{key}", c.getURL)
}

type SaveUrlRequest struct {
	URL string `json:"url"`
}

func (c URLController) saveURL(w http.ResponseWriter, r *http.Request) {
	var req SaveUrlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.New(w, ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}

	trimmedPrefix := strings.Trim(c.URLPrefix, "/")
	fmt.Println(trimmedPrefix, req.URL)
	if strings.HasPrefix(req.URL, trimmedPrefix) {
		response.New(w, ErrUrlNotAllowed.Error(), http.StatusBadRequest)
		return
	}

	if _, err := url.ParseRequestURI(req.URL); err != nil {
		response.New(w, ErrInvalidUrl.Error(), http.StatusBadRequest)
		return
	}

	key, err := c.URLService.SaveUrl(req.URL)
	if err != nil {
		response.New(w, ErrObfuscation.Error(), http.StatusInternalServerError)
		return
	}

	response.New(w, c.URLPrefix+key, http.StatusOK)
}

func (c URLController) getURL(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	pattern := regexp.MustCompile(fmt.Sprintf("[A-z0-9]{%d}", c.URLService.IndexService.KeyLimit))

	if !pattern.MatchString(key) {
		response.New(w, ErrInvalidKey.Error(), http.StatusBadRequest)
		return
	}

	url, err := c.URLService.GetURLByKey(key)
	if errors.Is(err, service.ErrNotFound) {
		response.New(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		response.New(w, ErrObfuscation.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, url, http.StatusPermanentRedirect)
}
