package middlewares

import (
	"github.com/isaacwallace123/GoUtils/color"
	"github.com/isaacwallace123/GoUtils/logger"
	"github.com/isaacwallace123/GoWeb/app/types"
)

type LoggingConfig struct {
	Enabled bool
}

// Logs before the handler runs (no response yet)
var LoggingPre = types.NewMiddlewareBuilder("logging_pre", &LoggingConfig{
	Enabled: true,
}, func(ctx *types.MiddlewareContext, cfg *LoggingConfig) error {
	if cfg.Enabled {
		methodColored := color.HTTPMethodToColor[ctx.Request.Method] + ctx.Request.Method + color.Reset
		logger.Info("%s %s", methodColored, ctx.Request.URL.Path)
	}
	return ctx.Next()
})

// Logs after the handler, when ResponseEntity exists
var LoggingPost = types.NewMiddlewareBuilder("logging_post", &LoggingConfig{
	Enabled: true,
}, func(ctx *types.MiddlewareContext, cfg *LoggingConfig) error {
	if cfg.Enabled && ctx.ResponseEntity != nil {
		methodColored := color.HTTPMethodToColor[ctx.Request.Method] + ctx.Request.Method + color.Reset
		logger.Info("%s %s %d", methodColored, ctx.Request.URL.Path, ctx.ResponseEntity.StatusCode)
	}
	return ctx.Next()
})
