package core

import (
	"log"
	"net/http"
	"reflect"

	"github.com/isaacwallace123/GoWeb/decorators"
	"github.com/isaacwallace123/GoWeb/httpstatus"
	"github.com/isaacwallace123/GoWeb/response"
)

type SpringRouter struct {
	mux *http.ServeMux
}

func NewRouter() *SpringRouter {
	return &SpringRouter{
		mux: http.NewServeMux(),
	}
}

func (r *SpringRouter) RegisterControllers() {
	routeTable := make(map[string]map[string]routeHandler)

	for _, route := range decorators.RegisteredRoutes() {
		for _, ctrl := range decorators.RegisteredControllers() {
			fullPath := ctrl.BasePath + route.Path

			if routeTable[fullPath] == nil {
				routeTable[fullPath] = make(map[string]routeHandler)
			}

			if _, exists := routeTable[fullPath][route.Method]; exists {
				log.Fatalf("Conflict: %s [%s] already registered", fullPath, route.Method)
			}

			routeTable[fullPath][route.Method] = routeHandler{
				Controller: ctrl.Instance,
				Method:     reflect.ValueOf(route.HandlerFunc),
			}
		}
	}

	for path, methodMap := range routeTable {
		r.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
			handlerEntry, ok := methodMap[req.Method]
			if !ok {
				response.Status(httpstatus.METHOD_NOT_ALLOWED).
					Body(map[string]string{"error": "Method Not Allowed"}).
					Send(w)
				return
			}

			handler := handlerEntry.Method
			handlerType := handler.Type()
			args := []reflect.Value{reflect.ValueOf(handlerEntry.Controller)}

			// Bind @RequestBody if applicable
			if handlerType.NumIn() == 2 {
				argType := handlerType.In(1)
				argVal, err := BindJSONBody(req, argType)
				if err != nil {
					response.Status(httpstatus.BAD_REQUEST).
						Body(map[string]string{"error": "Invalid JSON"}).
						Send(w)
					return
				}
				args = append(args, reflect.ValueOf(argVal))
			}

			results := handler.Call(args)

			if len(results) != 1 {
				response.Status(httpstatus.INTERNAL_SERVER_ERR).
					Body(map[string]string{"error": "Handler must return exactly one value"}).
					Send(w)
				return
			}

			resp, ok := results[0].Interface().(*response.ResponseEntity)
			if !ok {
				response.Status(httpstatus.INTERNAL_SERVER_ERR).
					Body(map[string]string{"error": "Return type must be *response.ResponseEntity"}).
					Send(w)
				return
			}

			resp.Send(w)
		})
	}
}

func (r *SpringRouter) Listen(addr string) error {
	return http.ListenAndServe(addr, r.mux)
}

type routeHandler struct {
	Controller any
	Method     reflect.Value
}
