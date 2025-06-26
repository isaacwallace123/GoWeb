# go-responseentity

A small Go package that mimics Java Spring Boot's `ResponseEntity` for clean, chainable HTTP responses.

## Example Usage

```go
import (
	"net/http"
	"github.com/yourusername/go-responseentity/response"
	"github.com/yourusername/go-responseentity/httpstatus"
)

func handler(w http.ResponseWriter, r *http.Request) {
	response.Status(httpstatus.CREATED).
		Body(map[string]string{"test": "test"}).
		Header("key", "value").
		Send(w)
}
