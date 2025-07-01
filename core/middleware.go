package core

import (
	"github.com/isaacwallace123/GoUtils/logger"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

var middlewares []Middleware

func Use(mw ...Middleware) {
	middlewares = append(middlewares, mw...)
}

func chainMiddleware(final http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		final = middlewares[i](final)
	}
	return final
}

var methodColors = map[string]string{
	"GET":    "\033[32m",
	"POST":   "\033[34m",
	"PUT":    "\033[33m",
	"DELETE": "\033[31m",
	"PATCH":  "\033[35m",
}

const reset = "\033[0m"

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		color, ok := methodColors[r.Method]
		if !ok {
			color = "\033[90m"
		}

		methodColored := color + r.Method + reset
		logger.Info("%s %s", methodColored, r.URL.Path)

		next.ServeHTTP(w, r)
	})
}
