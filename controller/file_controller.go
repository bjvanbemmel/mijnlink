package controller

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/bjvanbemmel/mijnlink/response"
	"github.com/bjvanbemmel/mijnlink/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type FileController struct {
	URLPrefix       string
	FileService     service.FileService
	UploadSizeLimit int
}

func (c FileController) InitRoutes(r *chi.Mux) {
	group := r.Group(nil)
	group.Use(
		middleware.AllowContentType("multipart/form-data"),
	)
	group.Post("/file", c.saveFile)
	group.Get("/file/{key}", c.getFile)
}

func (c FileController) saveFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm((1024 ^ 3) * int64(c.UploadSizeLimit))

	file, _, err := r.FormFile("file")
	if err != nil {
		response.New(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := c.FileService.SaveFile(file)
	if err != nil {
		response.New(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.New(w, c.URLPrefix+res, http.StatusOK)
}

func (c FileController) getFile(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	pattern := regexp.MustCompile(fmt.Sprintf("[A-z0-9]{%d}", c.FileService.IndexService.KeyLimit))

	if !pattern.MatchString(key) {
		response.New(w, ErrInvalidKey.Error(), http.StatusBadRequest)
		return
	}

	contents, err := c.FileService.GetFileByKey(key)
	if errors.Is(err, service.ErrNotFound) {
		response.New(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		response.New(w, ErrObfuscation.Error(), http.StatusInternalServerError)
		return
	}

	response.NewFile(w, []byte(contents))
}
