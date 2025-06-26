package core

import (
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/isaacwallace123/GoWeb/httpstatus"
	"github.com/isaacwallace123/GoWeb/response"
)

type Controller interface {
	Path() string
}

type routeEntry struct {
	Method     string
	Regex      *regexp.Regexp
	ParamNames []string
	Handler    reflect.Value
	CtrlValue  reflect.Value
}

type SpringRouter struct {
	routes []routeEntry
}

func NewRouter() *SpringRouter {
	return &SpringRouter{}
}

func RegisterControllers(instances ...Controller) *SpringRouter {
	r := NewRouter()

	for _, inst := range instances {
		path := inst.Path()
		re, paramNames := compilePathPattern(path)
		val := reflect.ValueOf(inst)
		typ := reflect.TypeOf(inst)

		for i := 0; i < typ.NumMethod(); i++ {
			m := typ.Method(i)
			methodName := strings.ToUpper(m.Name)

			if isHTTPMethod(methodName) {
				r.routes = append(r.routes, routeEntry{
					Method:     methodName,
					Regex:      re,
					ParamNames: paramNames,
					Handler:    val.Method(i),
					CtrlValue:  val,
				})
			}
		}
	}

	return r
}

func (r *SpringRouter) Listen(addr string) error {
	http.HandleFunc("/", r.dispatch)
	return http.ListenAndServe(addr, nil)
}

func (r *SpringRouter) dispatch(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if req.Method != route.Method {
			continue
		}

		matches := route.Regex.FindStringSubmatch(req.URL.Path)
		if matches == nil {
			continue
		}

		pathVars := map[string]string{}
		for i, name := range route.ParamNames {
			pathVars[name] = matches[i+1]
		}

		handlerType := route.Handler.Type()
		paramTypes := make([]reflect.Type, handlerType.NumIn())
		for i := 0; i < handlerType.NumIn(); i++ {
			paramTypes[i] = handlerType.In(i)
		}

		args, err := BindArguments(req, req.Context(), paramTypes, pathVars)
		if err != nil {
			response.Status(httpstatus.BAD_REQUEST).
				Body(map[string]string{"error": err.Error()}).
				Send(w)
			return
		}

		result := route.Handler.Call(args)

		if len(result) != 1 {
			response.Status(httpstatus.INTERNAL_SERVER_ERR).
				Body(map[string]string{"error": "Expected 1 return value"}).
				Send(w)
			return
		}

		resp, ok := result[0].Interface().(*response.ResponseEntity)
		if !ok {
			response.Status(httpstatus.INTERNAL_SERVER_ERR).
				Body(map[string]string{"error": "Invalid return type"}).
				Send(w)
			return
		}

		final := chainMiddleware(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			resp.Send(w)
		}))

		final.ServeHTTP(w, req)

		return
	}

	response.Status(httpstatus.NOT_FOUND).
		Body(map[string]string{"error": "Route not found"}).
		Send(w)
}

func compilePathPattern(path string) (*regexp.Regexp, []string) {
	paramRegex := regexp.MustCompile(`\{([^}]+)\}`)
	paramNames := []string{}
	regexStr := paramRegex.ReplaceAllStringFunc(path, func(m string) string {
		name := m[1 : len(m)-1]
		paramNames = append(paramNames, name)
		return "([^/]+)"
	})
	return regexp.MustCompile("^" + regexStr + "$"), paramNames
}

func isHTTPMethod(name string) bool {
	switch name {
	case "GET", "POST", "PUT", "DELETE":
		return true
	}
	return false
}
