package response

type ResponseEntity struct {
	status int
	header map[string]string
	body   any
}

func Status(code int) *ResponseEntity {
	return &ResponseEntity{
		status: code,
		header: make(map[string]string),
	}
}

func (r *ResponseEntity) Body(body any) *ResponseEntity {
	r.body = body
	return r
}

func (r *ResponseEntity) Header(key, value string) *ResponseEntity {
	r.header[key] = value
	return r
}
