package core

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
)

func BindArguments(
	req *http.Request,
	ctx context.Context,
	paramTypes []reflect.Type,
	pathVars map[string]string,
) ([]reflect.Value, error) {
	args := []reflect.Value{}
	ctx = WithPathVars(ctx, pathVars)
	ctx = WithQueryParams(ctx, req)
	ctx = WithHeaderMap(ctx, req.Header)

	args = append(args, reflect.ValueOf(ctx))

	for i := 1; i < len(paramTypes); i++ {
		t := paramTypes[i]

		// Try body (assume struct and method is POST/PUT)
		if t.Kind() == reflect.Struct && (req.Method == http.MethodPost || req.Method == http.MethodPut) {
			ptr := reflect.New(t).Interface()
			err := json.NewDecoder(req.Body).Decode(ptr)
			if err != nil {
				return nil, err
			}
			args = append(args, reflect.ValueOf(reflect.ValueOf(ptr).Elem().Interface()))
			continue
		}

		name := t.Name()

		if val := PathVar(ctx, name); val != "" {
			args = append(args, reflect.ValueOf(val))
			continue
		}

		// Try query param
		if val := QueryParam(ctx, name); val != "" {
			args = append(args, reflect.ValueOf(val))
			continue
		}

		// Try header
		if val := Header(ctx, name); val != "" {
			args = append(args, reflect.ValueOf(val))
			continue
		}

		// Default to empty string if nothing found
		args = append(args, reflect.Zero(t))
	}

	return args, nil
}
