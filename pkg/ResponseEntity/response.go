package ResponseEntity

import (
	"github.com/isaacwallace123/GoUtils/jsonutil"
	"net/http"
)

type ResponseEntity struct {
	StatusCode int
	Headers    map[string]string
	BodyData   any
}

// Build creates a new empty ResponseEntity.
func Build() *ResponseEntity {
	return &ResponseEntity{
		Headers: make(map[string]string),
	}
}

// Status starts a new response with an HTTP status and uses Build for consistency.
func Status(code int) *ResponseEntity {
	response := Build()
	response.StatusCode = code

	return response
}

// Body Chainable method to set body
func (response *ResponseEntity) Body(data any) *ResponseEntity {
	response.BodyData = data

	return response
}

// Header Chainable method to set headers
func (response *ResponseEntity) Header(key, value string) *ResponseEntity {
	response.Headers[key] = value

	return response
}

// Send writes the response
func (response *ResponseEntity) Send(writer http.ResponseWriter) {
	for k, v := range response.Headers {
		writer.Header().Set(k, v)
	}

	if response.BodyData != nil && response.StatusCode != http.StatusNoContent {
		writer.Header().Set("Content-Type", "application/json")
	}

	writer.WriteHeader(response.StatusCode)

	if response.BodyData != nil && response.StatusCode != http.StatusNoContent {
		jsonStr := jsonutil.ToString(response.BodyData)
		_, _ = writer.Write([]byte(jsonStr))
	}
}
