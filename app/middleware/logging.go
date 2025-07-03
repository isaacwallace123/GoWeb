package middleware

import (
	"github.com/isaacwallace123/GoUtils/logger"
	"net/http"
)

// Color codes
var methodColors = map[string]string{
	"GET":    "\033[32m",
	"POST":   "\033[34m",
	"PUT":    "\033[33m",
	"DELETE": "\033[31m",
	"PATCH":  "\033[35m",
}

const reset = "\033[0m"

// Pre-middleware: logs incoming request
func LoggingPre(ctx *PreMiddlewareContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		color := methodColors[r.Method]
		methodColored := color + r.Method + reset
		logger.Info("%s %s", methodColored, r.URL.Path)
		ctx.Handler.ServeHTTP(w, r)
	})
}

// Post-middleware: logs request + response status
func LoggingPost(ctx *PostMiddlewareContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		color := methodColors[r.Method]
		methodColored := color + r.Method + reset
		status := 0
		if ctx.Response != nil {
			status = ctx.Response.StatusCode
		}
		logger.Info("%s %s %d", methodColored, r.URL.Path, status)
		ctx.Handler.ServeHTTP(w, r)
	})
}
