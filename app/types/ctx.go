package types

import (
	"context"
	"net/http"
)

// contextKey is a custom type used to avoid key collisions in context values.
type contextKey string

// Keys used for storing values in context.
const (
	pathVarsKey    contextKey = "pathVars"    // For path parameters (e.g., /users/{id})
	queryParamsKey contextKey = "queryParams" // For URL query parameters (?foo=bar)
	headerMapKey   contextKey = "headerMap"   // For HTTP request headers
)

// WithPathVars adds path parameters to the context.
func WithPathVars(ctx context.Context, vars map[string]string) context.Context {
	return context.WithValue(ctx, pathVarsKey, vars)
}

// WithQueryParams adds query parameters from the request to the context.
func WithQueryParams(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, queryParamsKey, req.URL.Query())
}

// WithHeaderMap adds HTTP headers to the context.
func WithHeaderMap(ctx context.Context, headers http.Header) context.Context {
	return context.WithValue(ctx, headerMapKey, headers)
}

// PathVar retrieves a path parameter from the context by name.
// Returns an empty string if not found.
func PathVar(ctx context.Context, name string) string {
	if vars, ok := ctx.Value(pathVarsKey).(map[string]string); ok {
		return vars[name]
	}
	return ""
}

// QueryParam retrieves a single query parameter by name from the context.
// Returns an empty string if not found or if the parameter has no value.
func QueryParam(ctx context.Context, name string) string {
	if values, ok := ctx.Value(queryParamsKey).(map[string][]string); ok {
		if val, exists := values[name]; exists && len(val) > 0 {
			return val[0]
		}
	}
	return ""
}

// Header retrieves a specific header value from the context by name.
// Returns an empty string if not found.
func Header(ctx context.Context, name string) string {
	if headers, ok := ctx.Value(headerMapKey).(http.Header); ok {
		return headers.Get(name)
	}
	return ""
}
