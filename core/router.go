package core

import (
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
	for _, route := range decorators.RegisteredRoutes() {
		for _, ctrl := range decorators.RegisteredControllers() {
			fullPath := ctrl.BasePath + route.Path

			// Bind method like (*Controller).Create
			methodVal := reflect.ValueOf(route.HandlerFunc)
			methodType := methodVal.Type()

			r.mux.HandleFunc(fullPath, func(w http.ResponseWriter, req *http.Request) {
				if req.Method != route.Method {
					response.Status(httpstatus.METHOD_NOT_ALLOWED).
						Body(map[string]string{"error": "Method Not Allowed"}).
						Send(w)
					return
				}

				// Prepare arguments
				args := []reflect.Value{reflect.ValueOf(ctrl.Instance)}

				// If second argument exists, bind request body
				if methodType.NumIn() == 2 {
					argType := methodType.In(1)
					argVal, err := BindJSONBody(req, argType)
					if err != nil {
						response.Status(httpstatus.BAD_REQUEST).
							Body(map[string]string{"error": "Invalid JSON"}).
							Send(w)
						return
					}
					args = append(args, reflect.ValueOf(argVal))
				}

				// Call the method
				results := methodVal.Call(args)

				// Expecting: func(...) *response.ResponseEntity
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
}

func (r *SpringRouter) Listen(addr string) error {
	return http.ListenAndServe(addr, r.mux)
}
