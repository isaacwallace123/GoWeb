package middlewares

import (
	"github.com/isaacwallace123/GoUtils/logger"
	"github.com/isaacwallace123/GoWeb/app/types"
)

var methodColors = map[string]string{
	"GET":    "\033[32m",
	"POST":   "\033[34m",
	"PUT":    "\033[33m",
	"DELETE": "\033[31m",
	"PATCH":  "\033[35m",
}

const reset = "\033[0m"

// Pre-middleware: logs before the handler runs (no response yet)
func LoggingPre(ctx *types.MiddlewareContext) error {
	color := methodColors[ctx.Request.Method]
	methodColored := color + ctx.Request.Method + reset

	logger.Info("%s %s", methodColored, ctx.Request.URL.Path)

	return ctx.Next()
}

// Post-middleware: logs after the handler, when ResponseEntity exists
func LoggingPost(ctx *types.MiddlewareContext) error {
	color := methodColors[ctx.Request.Method]
	methodColored := color + ctx.Request.Method + reset

	logger.Info("%s %s %d", methodColored, ctx.Request.URL.Path, ctx.ResponseEntity.StatusCode)

	return ctx.Next()
}
