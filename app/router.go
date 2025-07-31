package app

import (
	"github.com/isaacwallace123/GoWeb/app/internal"
	"github.com/isaacwallace123/GoWeb/app/types"
	"net/http"
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
func (r *Router) Listen(addr string) error {
	return internal.ListenImpl(r.routes, addr)
}

// ServeHTTP allows Router to implement http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, handler := range r.resources {
		if handler(w, req) {
			return
		}
	}
	internal.Dispatch(r.routes, w, req)
}
