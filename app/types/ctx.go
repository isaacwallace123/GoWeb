package types

import (
	"context"
	"net/http"
)

type contextKey string

const (
	pathVarsKey    contextKey = "pathVars"
	queryParamsKey contextKey = "queryParams"
	headerMapKey   contextKey = "headerMap"
)

func WithPathVars(ctx context.Context, vars map[string]string) context.Context {
	return context.WithValue(ctx, pathVarsKey, vars)
}

func WithQueryParams(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, queryParamsKey, req.URL.Query())
}

func WithHeaderMap(ctx context.Context, headers http.Header) context.Context {
	return context.WithValue(ctx, headerMapKey, headers)
}

func PathVar(ctx context.Context, name string) string {
	if vars, ok := ctx.Value(pathVarsKey).(map[string]string); ok {
		return vars[name]
	}
	return ""
}

func QueryParam(ctx context.Context, name string) string {
	if values, ok := ctx.Value(queryParamsKey).(map[string][]string); ok {
		if val, exists := values[name]; exists && len(val) > 0 {
			return val[0]
		}
	}
	return ""
}

func Header(ctx context.Context, name string) string {
	if headers, ok := ctx.Value(headerMapKey).(http.Header); ok {
		return headers.Get(name)
	}
	return ""
}
