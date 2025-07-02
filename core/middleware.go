package core

import (
	"net/http"

	"github.com/isaacwallace123/GoUtils/logger"
	"github.com/isaacwallace123/GoWeb/response"
)
type Middleware func(next http.Handler, args ...any) http.Handler

const (
	PRE_MIDDLEWARE 	= 0
	POST_MIDDLEWARE = 1
)

var premiddlewares []Middleware
var postmiddlewares []Middleware

func Use(priority int, mw ...Middleware){
	if(priority == 0){
		premiddlewares = append(premiddlewares, mw...)
	}else{
		postmiddlewares = append(postmiddlewares, mw...)
	}
}

var methodColors = map[string]string{
	"GET":    "\033[32m",
	"POST":   "\033[34m",
	"PUT":    "\033[33m",
	"DELETE": "\033[31m",
	"PATCH":  "\033[35m",
}
const reset = "\033[0m"

func LoggingMiddleware(next http.Handler, args ...any) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		color, ok := methodColors[r.Method]
		if !ok {
			color = "\033[90m"
		}

		arg, ok := args[0].(response.ResponseEntity)

		if(!ok){
			logger.Info("cool")
		}

		methodColored := color + r.Method + reset
		logger.Info("%s %s %d", methodColored, r.URL.Path, arg.StatusCode)

		next.ServeHTTP(w, r)
	})
}