package app

import (
	"github.com/isaacwallace123/GoUtils/logger"
	"github.com/isaacwallace123/GoWeb/app/internal"
	"github.com/isaacwallace123/GoWeb/app/types"
	"net/http"
	"strings"
)

type Router struct {
	routes    []internal.CompiledRoute
	resources []func(http.ResponseWriter, *http.Request) bool
}

// NewRouter creates a new Router.
func NewRouter() *Router {
	return &Router{}
}

// Make RegisterControllers a method on *Router
func (r *Router) RegisterControllers(controllers ...types.Controller) {
	r.routes = internal.RegisterControllersImpl(controllers...)
}

// Listen starts the HTTP server.
func (r *Router) Listen(addr string) error { return internal.ListenImpl(r, addr) }

// ServeHTTP first tries static handlers, then dispatches dynamic routes.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, handler := range r.resources {
		if handler(w, req) {
			return
		}
	}

	internal.Dispatch(r.routes, w, req)
}

func (r *Router) ListAllRoutes() []types.Route {
	routes := make([]types.Route, 0, len(r.routes))

	for _, cr := range r.routes {
		routes = append(routes, types.Route{
			Method:  cr.Method,
			Path:    cr.Path,
			Handler: cr.HandlerName,
		})
	}

	return routes
}

// UseStatic registers a static file handler for the given URL prefix and directory.
func (r *Router) UseStatic(prefix, dir string) {
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	if len(prefix) > 1 && strings.HasSuffix(prefix, "/") {
		prefix = strings.TrimSuffix(prefix, "/")
	}

	fs := http.FileServer(http.Dir(dir))
	handler := func(w http.ResponseWriter, req *http.Request) bool {
		path := req.URL.Path

		if path == prefix {
			http.Redirect(w, req, prefix+"/", http.StatusMovedPermanently)
			logger.Info("[Static] Redirected: %s → %s/", path, prefix)
			return true
		}

		if strings.HasPrefix(path, prefix+"/") {
			logger.Info("[Static] %s → %s (%s)", prefix, dir, path)
			http.StripPrefix(prefix, fs).ServeHTTP(w, req)
			return true
		}
		return false
	}

	r.resources = append(r.resources, handler)
	logger.Info("[Static] Registered: %-12s → %s", prefix, dir)
}
