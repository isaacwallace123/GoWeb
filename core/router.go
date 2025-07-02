package core

import (
	"context"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/isaacwallace123/GoWeb/response"
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

		result := route.Handler.Call(args)
		if len(result) != 1 {
			InternalServerException("Expected 1 return value").Send(w)
			return
		}

		resp, ok := result[0].Interface().(*response.ResponseEntity)
		if !ok {
			InternalServerException("Invalid return type").Send(w)
			return
		}

		var final http.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {})

		// Pre-middlewares
		for _, ware  := range premiddlewares {
			final = ware(final)
		}

		resp.Send(w)

		// Post-middlewares
		for _, ware  := range postmiddlewares {
			final = ware(final, resp)
		}

		final.ServeHTTP(w, req)
		return
	}

	NotFoundException("Route not found").Send(w)
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
