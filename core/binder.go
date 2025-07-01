// core/binder.go
package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

func BindArguments(
	req *http.Request,
	ctx context.Context,
	paramTypes []reflect.Type,
	pathVars map[string]string,
	argNames []string,
) ([]reflect.Value, error) {
	args := []reflect.Value{}
	start := 0

	hasCtx := len(paramTypes) > 0 && paramTypes[0] == reflect.TypeOf((*context.Context)(nil)).Elem()
	if hasCtx {
		args = append(args, reflect.ValueOf(ctx))
		start = 1
	}

	ctx = WithPathVars(ctx, pathVars)
	ctx = WithQueryParams(ctx, req)
	ctx = WithHeaderMap(ctx, req.Header)

	for i := start; i < len(paramTypes); i++ {
		t := paramTypes[i]
		argIdx := i - start

		name := ""
		if argIdx < len(argNames) {
			name = argNames[argIdx]
		}

		if t.Kind() == reflect.Struct && (req.Method == http.MethodPost || req.Method == http.MethodPut) {
			ptr := reflect.New(t).Interface()
			err := json.NewDecoder(req.Body).Decode(ptr)
			if err != nil {
				return nil, err
			}
			args = append(args, reflect.ValueOf(reflect.ValueOf(ptr).Elem().Interface()))
			continue
		}

		if val, ok := pathVars[name]; ok {
			switch t.Kind() {
			case reflect.Int:
				intVal, err := strconv.Atoi(val)
				if err != nil {
					return nil, fmt.Errorf("invalid int for %s: %v", name, err)
				}
				args = append(args, reflect.ValueOf(intVal))
			default:
				args = append(args, reflect.ValueOf(val))
			}
			continue
		}

		args = append(args, reflect.Zero(t))
	}

	return args, nil
}
