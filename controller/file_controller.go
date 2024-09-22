package controller

import (
	"net/http"

	"github.com/bjvanbemmel/mijnlink/response"
	"github.com/bjvanbemmel/mijnlink/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type FileController struct {
	FileService     service.FileService
	UploadSizeLimit int
}

func (c FileController) InitRoutes(r *chi.Mux) {
	group := r.Group(nil)
	group.Use(
		middleware.AllowContentType("multipart/form-data"),
	)
	group.Post("/file", c.saveFile)
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

	response.New(w, res, http.StatusOK)
}
