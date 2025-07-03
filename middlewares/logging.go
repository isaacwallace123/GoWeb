package middlewares

import (
	"github.com/isaacwallace123/GoUtils/color"
	"github.com/isaacwallace123/GoUtils/logger"
	"github.com/isaacwallace123/GoWeb/app/types"
)

// Pre-middleware: logs before the handler runs (no response yet)
func LoggingPre(ctx *types.MiddlewareContext) error {
	methodColored := color.HTTPMethodToColor[ctx.Request.Method] + ctx.Request.Method + color.Reset

	logger.Info("%s %s", methodColored, ctx.Request.URL.Path)

	return ctx.Next()
}

// Post-middleware: logs after the handler, when ResponseEntity exists
func LoggingPost(ctx *types.MiddlewareContext) error {
	methodColored := color.HTTPMethodToColor[ctx.Request.Method] + ctx.Request.Method + color.Reset

	logger.Info("%s %s %d", methodColored, ctx.Request.URL.Path, ctx.ResponseEntity.StatusCode)

	return ctx.Next()
}
