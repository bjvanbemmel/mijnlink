package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Result struct {
	Value string `json:"value"`
}

func (r Result) JSON() []byte {
	raw, _ := json.Marshal(r)
	return raw
}

func New(w http.ResponseWriter, value string, status int) {
	res := Result{
		Value: value,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res.JSON())
}

func NewFile(w http.ResponseWriter, value []byte, filename string) {
	w.Header().Add("Content-Type", http.DetectContentType(value))
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Add("Content-Size", fmt.Sprintf("attachment; filename=%d", len(value)))
	w.Header().Add("Access-Control-Expose-Headers", "Content-Disposition,Content-Size")
	w.WriteHeader(http.StatusOK)
	w.Write(value)
}

func NewFileHeaders(w http.ResponseWriter, value []byte, filename string) {
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Add("Content-Size", fmt.Sprintf("attachment; filename=%d", len(value)))
	w.Header().Add("Access-Control-Expose-Headers", "Content-Disposition,Content-Size")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}
