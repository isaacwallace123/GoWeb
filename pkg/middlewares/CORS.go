package middlewares

import (
	"fmt"
	"github.com/isaacwallace123/GoUtils/color"
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

var CORS_TAG = fmt.Sprintf("%sCORS%s", color.BrightCyan, color.Reset)

var CORS = types.NewMiddlewareBuilder("cors", &CORSConfig{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Content-Type", "Authorization"},
	AllowCredentials: true,
}, func(ctx *types.MiddlewareContext, config *CORSConfig) error {
	origin := ctx.Request.Header.Get("Origin")

	// Check and set all headers if needed
	if origin != "" {
		if isOriginAllowed(origin, config.AllowedOrigins) {
			headers := ctx.ResponseWriter.Header()
			headers.Set("Access-Control-Allow-Origin", origin)
			headers.Set("Vary", "Origin")
			headers.Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
			headers.Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))

			if config.AllowCredentials {
				headers.Set("Access-Control-Allow-Credentials", "true")
			}

			logger.Debug("%s Origin allowed: %s", CORS_TAG, origin)
		} else {
			logger.Warn("%s Disallowed origin: %s", CORS_TAG, origin)
		}
	} else {
		logger.Debug("%s No Origin header present", CORS_TAG)
	}

	// Always intercept OPTIONS
	if ctx.Request.Method == http.MethodOptions {
		logger.Debug("%s Preflight OPTIONS request intercepted", CORS_TAG)

		ctx.ResponseWriter.WriteHeader(http.StatusNoContent)

		return nil
	}

	// Enforce method restriction on all other requests
	if !isMethodAllowed(ctx.Request.Method, config.AllowedMethods) {
		logger.Warn("%s Method not allowed: %s", CORS_TAG, color.HTTPMethodToColor[ctx.Request.Method]+ctx.Request.Method+color.Reset)

		ctx.ResponseWriter.WriteHeader(http.StatusMethodNotAllowed)

		return nil
	}

	return ctx.Next()
}).WithInit(func(config *CORSConfig) {
	logger.Info("%s Middleware successfully initialized", CORS_TAG)
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
