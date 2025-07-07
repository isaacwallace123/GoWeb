package app

import "github.com/isaacwallace123/GoWeb/app/types"

// Register pre-middleware (as Middleware interface, not just funcs)
func Use(mw ...types.Middleware) {
	types.PreMiddlewares = append(types.PreMiddlewares, mw...)
}

// Register post-middleware
func UseAfter(mw ...types.Middleware) {
	types.PostMiddlewares = append(types.PostMiddlewares, mw...)
}

// Optional accessors
func Pre() []types.MiddlewareFunc {
	return types.ConvertMiddewaresToFuncs(types.PreMiddlewares)
}

func Post() []types.MiddlewareFunc {
	return types.ConvertMiddewaresToFuncs(types.PostMiddlewares)
}
