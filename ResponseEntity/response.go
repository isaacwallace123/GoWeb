package ResponseEntity

import (
	"encoding/json"
	"net/http"
)

type ResponseEntity struct {
	StatusCode int
	Headers    map[string]string
	BodyData   any
}

// Start a new response with HTTP status
func Status(code int) *ResponseEntity {
	return &ResponseEntity{
		StatusCode: code,
		Headers:    make(map[string]string),
	}
}

// Chainable method to set body
func (r *ResponseEntity) Body(data any) *ResponseEntity {
	r.BodyData = data

	return r
}

// Chainable method to set headers
func (r *ResponseEntity) Header(key, value string) *ResponseEntity {
	r.Headers[key] = value

	return r
}

// Send writes the response
func (r *ResponseEntity) Send(w http.ResponseWriter) {
	for k, v := range r.Headers {
		w.Header().Set(k, v)
	}

	if r.BodyData != nil {
		w.Header().Set("Content-Type", "application/json")
	}

	w.WriteHeader(r.StatusCode)

	if r.BodyData != nil {
		_ = json.NewEncoder(w).Encode(r.BodyData)
	}
}
