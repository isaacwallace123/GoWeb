package types

import (
	ResponseEntity2 "github.com/isaacwallace123/GoWeb/ResponseEntity"
	"net/http"
)

type MiddlewareContext struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	ResponseEntity *ResponseEntity2.ResponseEntity
	Index          int
	Chain          []MiddlewareFunc
}

type MiddlewareFunc func(ctx *MiddlewareContext) error

func (ctx *MiddlewareContext) Next() error {
	ctx.Index++
	if ctx.Index < len(ctx.Chain) {
		return ctx.Chain[ctx.Index](ctx)
	}
	return nil
}

var PreMiddlewares []MiddlewareFunc
var PostMiddlewares []MiddlewareFunc
