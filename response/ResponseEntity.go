package response

import (
	"encoding/json"
	"net/http"
)

type ResponseEntity struct {
	status  int
	headers map[string]string
	body    any
}

func Status(status int) *ResponseEntity {
	return &ResponseEntity{
		status:  status,
		headers: make(map[string]string),
	}
}

func (r *ResponseEntity) Body(body any) *ResponseEntity {
	r.body = body
	return r
}

func (r *ResponseEntity) Header(key, value string) *ResponseEntity {
	r.headers[key] = value
	return r
}

func (r *ResponseEntity) Send(w http.ResponseWriter) {
	for k, v := range r.headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(r.status)
	if r.body != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(r.body)
	}
}
