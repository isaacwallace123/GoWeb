package middlewares

import (
	"github.com/isaacwallace123/GoUtils/logger"
	"net/http"
	"strings"

	"github.com/isaacwallace123/GoWeb/app/types"
)

type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

var CORS = types.NewMiddlewareBuilder("cors", &CORSConfig{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Content-Type", "Authorization"},
	AllowCredentials: true,
}, func(ctx *types.MiddlewareContext, config *CORSConfig) error {
	origin := ctx.Request.Header.Get("Origin")

	// Check and set all headers if needed
	if origin != "" && isOriginAllowed(origin, config.AllowedOrigins) {
		headers := ctx.ResponseWriter.Header()
		headers.Set("Access-Control-Allow-Origin", origin)
		headers.Set("Vary", "Origin")
		headers.Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
		headers.Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
		if config.AllowCredentials {
			headers.Set("Access-Control-Allow-Credentials", "true")
		}
	}

	// Always intercept OPTIONS
	if ctx.Request.Method == http.MethodOptions {
		ctx.ResponseWriter.WriteHeader(http.StatusNoContent)
		return nil
	}

	// Enforce method restriction on all other requests
	if !isMethodAllowed(ctx.Request.Method, config.AllowedMethods) {
		ctx.ResponseWriter.WriteHeader(http.StatusMethodNotAllowed)
		return nil
	}

	return ctx.Next()
}).WithInit(func(config *CORSConfig) {
	logger.Info("CORS Middleware successfully initialized")
})

func isOriginAllowed(origin string, allowed []string) bool {
	for _, o := range allowed {
		if o == "*" || o == origin {
			return true
		}
	}
	return false
}

func isMethodAllowed(method string, allowed []string) bool {
	for _, m := range allowed {
		if strings.EqualFold(m, method) {
			return true
		}
	}
	return false
}
