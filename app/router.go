package app

import (
	"context"
	"github.com/isaacwallace123/GoWeb/app/middleware"
	"github.com/isaacwallace123/GoWeb/exception"
	"github.com/isaacwallace123/GoWeb/response"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

type Router struct {
	routes []CompiledRoute
}

type CompiledRoute struct {
	Method     string
	Regex      *regexp.Regexp
	ParamNames []string
	Handler    reflect.Value
	CtrlValue  reflect.Value
}

type RouteEntry struct {
	Method  string // e.g. "GET", "POST"
	Path    string // e.g. "/", "/{userid}"
	Handler string // e.g. "Get", "Post", "GetAll"
}

type Controller interface {
	BasePath() string
	Routes() []RouteEntry
}

func NewRouter() *Router {
	return &Router{}
}

func RegisterControllers(controllers ...Controller) *Router {
	r := NewRouter()

	for _, ctrl := range controllers {
		val := reflect.ValueOf(ctrl)
		typ := reflect.TypeOf(ctrl)

		for _, entry := range ctrl.Routes() {
			fullPath := joinPath(ctrl.BasePath(), entry.Path)
			re, paramNames := compilePathPattern(fullPath)

			if _, ok := typ.MethodByName(entry.Handler); !ok {
				panic("Handler method not found: " + entry.Handler)
			}

			r.routes = append(r.routes, CompiledRoute{
				Method:     strings.ToUpper(entry.Method),
				Regex:      re,
				ParamNames: paramNames,
				Handler:    val.MethodByName(entry.Handler),
				CtrlValue:  val,
			})
		}
	}

	return r
}

func (r *Router) Listen(addr string) error {
	http.HandleFunc("/", r.dispatch)
	return http.ListenAndServe(addr, nil)
}

func (r *Router) dispatch(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if req.Method != route.Method {
			continue
		}

		normalizedPath := normalizePath(req.URL.Path)
		matches := route.Regex.FindStringSubmatch(normalizedPath)
		if matches == nil {
			continue
		}

		pathVars := extractPathVars(route.ParamNames, matches[1:])
		paramTypes := getParamTypes(route.Handler.Type())
		argNames := buildArgNames(paramTypes, route.ParamNames)

		args, err := BindArguments(req, req.Context(), paramTypes, pathVars, argNames)
		if err != nil {
			return
		}

		var responseEntity *response.ResponseEntity

		controllerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			result := route.Handler.Call(args)
			if len(result) != 1 {
				exception.InternalServerException("Expected 1 return value").Send(w)
				return
			}
			resp, ok := result[0].Interface().(*response.ResponseEntity)
			if !ok {
				exception.InternalServerException("Invalid return type").Send(w)
				return
			}
			responseEntity = resp
		})

		var handler http.Handler = controllerHandler
		for i := len(middleware.PreMiddlewares) - 1; i >= 0; i-- {
			handler = middleware.PreMiddlewares[i](&middleware.PreMiddlewareContext{
				Handler: handler,
				Request: req,
			})
		}

		recorder := &ResponseRecorder{ResponseWriter: w}
		handler.ServeHTTP(recorder, req)

		resp := responseEntity
		if resp == nil && recorder.Entity != nil {
			resp = recorder.Entity
		}

		var postHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if resp != nil {
				resp.Send(w)
			} else {
				exception.InternalServerException("No response entity").Send(w)
			}
		})

		for i := len(middleware.PostMiddlewares) - 1; i >= 0; i-- {
			postHandler = middleware.PostMiddlewares[i](&middleware.PostMiddlewareContext{
				Handler:  postHandler,
				Request:  req,
				Response: resp,
			})
		}

		postHandler.ServeHTTP(w, req)

		return
	}

	exception.NotFoundException("Route not found").Send(w)
}

type ResponseRecorder struct {
	http.ResponseWriter
	Entity *response.ResponseEntity
}

func (rr *ResponseRecorder) WriteHeader(statusCode int) {
	rr.ResponseWriter.WriteHeader(statusCode)
}

func (rr *ResponseRecorder) Write(b []byte) (int, error) {
	return rr.ResponseWriter.Write(b)
}

func normalizePath(path string) string {
	if path != "/" {
		return strings.TrimRight(path, "/")
	}
	return path
}

func joinPath(base, suffix string) string {
	base = strings.TrimRight(base, "/")
	suffix = strings.TrimLeft(suffix, "/")

	full := "/" + strings.TrimLeft(base+"/"+suffix, "/")

	if suffix == "" {
		full = base
	}

	return full
}

func compilePathPattern(path string) (*regexp.Regexp, []string) {
	paramRegex := regexp.MustCompile(`\{([^}]+)\}`)
	paramNames := []string{}

	regexStr := paramRegex.ReplaceAllStringFunc(path, func(m string) string {
		name := m[1 : len(m)-1]
		paramNames = append(paramNames, name)
		return "([^/]+)"
	})

	return regexp.MustCompile("^" + regexStr + "/?$"), paramNames
}

func extractPathVars(names, values []string) map[string]string {
	vars := make(map[string]string, len(names))
	for i := range names {
		vars[names[i]] = values[i]
	}
	return vars
}

func getParamTypes(handlerType reflect.Type) []reflect.Type {
	params := make([]reflect.Type, handlerType.NumIn())
	for i := 0; i < handlerType.NumIn(); i++ {
		params[i] = handlerType.In(i)
	}
	return params
}

func buildArgNames(paramTypes []reflect.Type, routeParams []string) []string {
	hasContext := len(paramTypes) > 0 && paramTypes[0] == reflect.TypeOf((*context.Context)(nil)).Elem()
	if hasContext {
		return append([]string{""}, routeParams...)
	}
	return routeParams
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.dispatch(w, req)
}
