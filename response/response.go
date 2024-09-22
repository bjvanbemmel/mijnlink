package response

import (
	"encoding/json"
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
