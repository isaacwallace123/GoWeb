package types

import (
	"net/http"
)

// MiddlewareContext carries request/response information and controls middleware flow.
type MiddlewareContext struct {
	Request        *http.Request       // Incoming HTTP request
	ResponseWriter http.ResponseWriter // Response writer for sending output
	ResponseEntity *ResponseEntity     // Optional response to be sent later
	Index          int                 // Current index in the middleware chain
	Chain          []MiddlewareFunc    // Ordered list of middleware to execute
}

// MiddlewareFunc represents a single middleware function.
// It receives a MiddlewareContext and returns an error if any.
type MiddlewareFunc func(ctx *MiddlewareContext) error

// Next advances the middleware chain to the next function.
// It returns any error produced by the next middleware.
func (ctx *MiddlewareContext) Next() error {
	ctx.Index++
	if ctx.Index < len(ctx.Chain) {
		return ctx.Chain[ctx.Index](ctx)
	}
	return nil // End of middleware chain
}

// PreMiddlewares holds globally registered middleware that runs before the handler.
var PreMiddlewares []MiddlewareFunc

// PostMiddlewares holds globally registered middleware that runs after the handler.
var PostMiddlewares []MiddlewareFunc
