package app

import "github.com/isaacwallace123/GoWeb/app/types"

// Register pre-middleware
func Use(mw ...types.MiddlewareFunc) {
	types.PreMiddlewares = append(types.PreMiddlewares, mw...)
}

// Register post-middleware
func UseAfter(mw ...types.MiddlewareFunc) {
	types.PostMiddlewares = append(types.PostMiddlewares, mw...)
}

func Pre() []types.MiddlewareFunc  { return types.PreMiddlewares }
func Post() []types.MiddlewareFunc { return types.PostMiddlewares }
