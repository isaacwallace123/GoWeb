package core

import "net/http"

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

// Example: Logging middleware
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println("[LOG]", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
