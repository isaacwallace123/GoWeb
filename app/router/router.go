package app

import (
	"encoding/json"
	"net/http"

	"github.com/isaacwallace123/GoWeb/response"
)

type HandlerFunc func(r *http.Request) *response.ResponseEntity

type Router struct {
	mux *http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

func (r *Router) Get(path string, handler HandlerFunc) {
	r.mux.HandleFunc(path, wrapHandler("GET", handler))
}

func (r *Router) Post(path string, handler HandlerFunc) {
	r.mux.HandleFunc(path, wrapHandler("POST", handler))
}

func (r *Router) Put(path string, handler HandlerFunc) {
	r.mux.HandleFunc(path, wrapHandler("PUT", handler))
}

func (r *Router) Del(path string, handler HandlerFunc) {
	r.mux.HandleFunc(path, wrapHandler("DELETE", handler))
}

func (r *Router) Listen(addr string) error {
	return http.ListenAndServe(addr, r.mux)
}

func wrapHandler(method string, handler HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		resp := handler(r)
		for k, v := range resp.header {
			w.Header().Set(k, v)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.status)
		if resp.Body != nil {
			json.NewEncoder(w).Encode(resp.Body)
		}
	}
}
