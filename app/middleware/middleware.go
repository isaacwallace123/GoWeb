package middleware

import (
	"github.com/isaacwallace123/GoWeb/response"
	"net/http"
)

type PreMiddleware func(*PreMiddlewareContext) http.Handler
type PostMiddleware func(*PostMiddlewareContext) http.Handler

type PreMiddlewareContext struct {
	Handler http.Handler
	Request *http.Request
}

type PostMiddlewareContext struct {
	Handler  http.Handler
	Request  *http.Request
	Response *response.ResponseEntity
}

var PreMiddlewares []PreMiddleware
var PostMiddlewares []PostMiddleware

func Use(mw ...PreMiddleware) {
	PreMiddlewares = append(PreMiddlewares, mw...)
}

func UseAfter(mw ...PostMiddleware) {
	PostMiddlewares = append(PostMiddlewares, mw...)
}
