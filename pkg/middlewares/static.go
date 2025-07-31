package middlewares

import (
	"github.com/isaacwallace123/GoUtils/logger"
	"github.com/isaacwallace123/GoWeb/app/types"
	"net/http"
	"strings"
)

var StaticMiddleware = types.NewMiddlewareBuilder("static", &types.StaticConfig{}, func(ctx *types.MiddlewareContext, cfg *types.StaticConfig) error {
	if cfg.Path == "" || cfg.Directory == "" {
		return ctx.Next()
	}
	if strings.HasPrefix(ctx.Request.URL.Path, cfg.Path) {
		fileServer := http.FileServer(http.Dir(cfg.Directory))
		http.StripPrefix(cfg.Path, fileServer).ServeHTTP(ctx.ResponseWriter, ctx.Request)
		return nil
	}
	return ctx.Next()
})

func StaticMiddlewares(statics []types.StaticConfig) []types.Middleware {
	var staticMiddlewares []types.Middleware

	for _, s := range statics {
		staticMiddlewares = append(staticMiddlewares, NewStaticMiddleware(s.Path, s.Directory))
	}

	return staticMiddlewares
}

func NewStaticMiddleware(path, directory string) types.Middleware {
	logger.Info("[Static] Registered: %-12s → %s", path, directory)

	return StaticMiddleware.WithInit(func(cfg *types.StaticConfig) {
		cfg.Path = path
		cfg.Directory = directory
	})
}
