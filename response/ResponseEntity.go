package response

type ResponseEntity struct {
	StatusCode int
	Headers    map[string]string
	body       any
}

func Status(code int) *ResponseEntity {
	return &ResponseEntity{
		StatusCode: code,
		Headers:    make(map[string]string),
	}
}

func (r *ResponseEntity) Body(body any) *ResponseEntity {
	r.body = body
	return r
}

func (r *ResponseEntity) Header(key, value string) *ResponseEntity {
	r.Headers[key] = value
	return r
}

func (r *ResponseEntity) GetBody() any {
	return r.body
}
