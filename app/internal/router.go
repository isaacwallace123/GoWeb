package internal

import (
	"context"
	"fmt"
	"github.com/isaacwallace123/GoWeb/pkg/exception"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/isaacwallace123/GoWeb/app/types"
)

// CompiledRoute struct remains unchanged
type CompiledRoute struct {
	Path        string
	Method      string
	Regex       *regexp.Regexp
	ParamNames  []string
	Handler     reflect.Value
	HandlerName string
	CtrlValue   reflect.Value
}

func RegisterControllersImpl(controllers ...types.Controller) []CompiledRoute {
	var compiled []CompiledRoute

	for _, ctrl := range controllers {
		val := reflect.ValueOf(ctrl)
		typ := reflect.TypeOf(ctrl)

		for _, entry := range ctrl.Routes() {
			fullPath := joinPath(ctrl.BasePath(), entry.Path)
			re, paramNames := compilePathPattern(fullPath)

			if _, ok := typ.MethodByName(entry.Handler); !ok {
				panic("Handler method not found: " + entry.Handler)
			}

			compiled = append(compiled, CompiledRoute{
				Path:        fullPath,
				Method:      strings.ToUpper(entry.Method),
				Regex:       re,
				ParamNames:  paramNames,
				Handler:     val.MethodByName(entry.Handler),
				HandlerName: fmt.Sprintf("%T.%s", ctrl, entry.Handler),
				CtrlValue:   val,
			})
		}
	}

	return compiled
}

func ListenImpl(router http.Handler, addr string) error {
	return http.ListenAndServe(addr, router)
}

func Dispatch(routes []CompiledRoute, w http.ResponseWriter, req *http.Request) {
	normalizedPath := normalizePath(req.URL.Path)

	for _, route := range routes {
		matches := route.Regex.FindStringSubmatch(normalizedPath)
		if matches == nil {
			continue
		}

		if req.Method != route.Method && req.Method != http.MethodOptions {
			continue
		}

		pathVars := extractPathVars(route.ParamNames, matches[1:])
		paramTypes := getParamTypes(route.Handler.Type())
		argNames := buildArgNames(paramTypes, route.ParamNames)
		args, err := BindArguments(req, req.Context(), paramTypes, pathVars, argNames)
		if err != nil {
			return
		}

		// --- Controller-level middleware
		var ctrlPre []types.Middleware
		if ctrl, ok := route.CtrlValue.Interface().(interface{ PreMiddleware() []types.Middleware }); ok {
			ctrlPre = ctrl.PreMiddleware()
		}

		var ctrlPost []types.Middleware
		if ctrl, ok := route.CtrlValue.Interface().(interface{ PostMiddleware() []types.Middleware }); ok {
			ctrlPost = ctrl.PostMiddleware()
		}

		// --- Build the chain
		chain := make([]types.MiddlewareFunc, 0,
			len(types.PreMiddlewares)+len(ctrlPre)+1+len(ctrlPost)+len(types.PostMiddlewares),
		)

		chain = append(chain, types.ConvertMiddewaresToFuncs(types.PreMiddlewares)...)
		chain = append(chain, types.ConvertMiddewaresToFuncs(ctrlPre)...)

		if req.Method != http.MethodOptions {
			chain = append(chain, func(ctx *types.MiddlewareContext) error {
				result := route.Handler.Call(args)
				if len(result) != 1 {
					exception.InternalServerException("Expected 1 return value").Send(w)
					return nil
				}
				if resp, ok := result[0].Interface().(*types.ResponseEntity); ok {
					ctx.ResponseEntity = resp
				}
				return ctx.Next()
			})
		}

		chain = append(chain, types.ConvertMiddewaresToFuncs(ctrlPost)...)
		chain = append(chain, types.ConvertMiddewaresToFuncs(types.PostMiddlewares)...)

		mwCtx := &types.MiddlewareContext{
			Request:        req,
			ResponseWriter: w,
			ResponseEntity: nil,
			Index:          -1,
			Chain:          chain,
		}

		_ = mwCtx.Next()

		if mwCtx.ResponseEntity != nil {
			mwCtx.ResponseEntity.Send(w)
		}
		return
	}

	if req.Method == http.MethodOptions {
		mwCtx := &types.MiddlewareContext{
			Request:        req,
			ResponseWriter: w,
			ResponseEntity: nil,
			Index:          -1,
			Chain:          types.ConvertMiddewaresToFuncs(types.PreMiddlewares),
		}
		_ = mwCtx.Next()
		return
	}

	// Fallback
	exception.NotFoundException("Route not found").Send(w)
}

// --- Helper functions (unchanged) ---

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
