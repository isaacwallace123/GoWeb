package types

import (
	"github.com/isaacwallace123/GoUtils/logger"
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
type MiddlewareFunc func(ctx *MiddlewareContext) error

// Middleware represents a middleware object with a Func method.
// This allows middleware structs to be registered and called.
type Middleware interface {
	Func() MiddlewareFunc
}

// Next advances the middleware chain to the next function.
// It returns any error produced by the next middleware.
func (ctx *MiddlewareContext) Next() error {
	ctx.Index++
	if ctx.Index < len(ctx.Chain) {
		return ctx.Chain[ctx.Index](ctx)
	}
	return nil // End of middleware chain
}

// PreMiddlewares holds globally registered middleware objects.
var PreMiddlewares []Middleware

// PostMiddlewares holds globally registered middleware objects.
var PostMiddlewares []Middleware

// --- Middleware Builder Pattern --- \\

// MiddlewareBuilder is a reusable, typed middleware object with attached config and logic.
type MiddlewareBuilder[T any] struct {
	Config         *T
	Handler        MiddlewareFunc
	OnErrorHandler func(ctx *MiddlewareContext, err error)
}

// Func allows the builder to be treated as a Middleware interface.
func (middleware *MiddlewareBuilder[T]) Func() MiddlewareFunc {
	return func(ctx *MiddlewareContext) error {
		err := middleware.Handler(ctx)

		if err != nil && middleware.OnErrorHandler != nil {
			middleware.OnErrorHandler(ctx, err)

			return nil
		}

		return err
	}
}

// WithInit is a script that the middleware's creator can tune to make it run once after the middleware is hooked
func (middleware *MiddlewareBuilder[T]) WithInit(initFn func(*T)) *MiddlewareBuilder[T] {
	initFn(middleware.Config)
	return middleware
}

func (middleware *MiddlewareBuilder[T]) OnError(handler func(ctx *MiddlewareContext, err error)) *MiddlewareBuilder[T] {
	middleware.OnErrorHandler = handler
	return middleware
}

// NewMiddlewareBuilder constructs a typed middleware with a config and handler function.
func NewMiddlewareBuilder[T any](
	name string,
	defaultConfig *T,
	handler func(ctx *MiddlewareContext, config *T) error,
) *MiddlewareBuilder[T] {
	builder := &MiddlewareBuilder[T]{
		Config: defaultConfig,
		Handler: func(ctx *MiddlewareContext) error {
			return handler(ctx, defaultConfig)
		},
		OnErrorHandler: func(ctx *MiddlewareContext, err error) {
			logger.Error("Middleware '%s' error: %v\n", name, err)
		},
	}
	return builder
}

func ConvertMiddewaresToFuncs(mw []Middleware) []MiddlewareFunc {
	funcs := make([]MiddlewareFunc, len(mw))
	for i, m := range mw {
		funcs[i] = m.Func()
	}
	return funcs
}
