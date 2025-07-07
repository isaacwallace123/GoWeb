package HttpMethod

const (
	GET     string = "GET"
	POST    string = "POST"
	PUT     string = "PUT"
	DELETE  string = "DELETE"
	PATCH   string = "PATCH"
	HEAD    string = "HEAD"
	OPTIONS string = "OPTIONS"
	CONNECT string = "CONNECT"
	TRACE   string = "TRACE"
)

var allMethods = []string{
	GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS, CONNECT, TRACE,
}

func IsValid(method string) bool {
	for _, m := range allMethods {
		if string(m) == method {
			return true
		}
	}
	return false
}
