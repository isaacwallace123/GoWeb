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
	argNames []string,
) ([]reflect.Value, error) {
	args := []reflect.Value{}

	ctx = WithPathVars(ctx, pathVars)
	ctx = WithQueryParams(ctx, req)
	ctx = WithHeaderMap(ctx, req.Header)

	//args = append(args, reflect.ValueOf(ctx))

	for i := 0; i < len(paramTypes); i++ {
		t := paramTypes[i]

		name := ""

		if i < len(argNames) {
			name = argNames[i]
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

		if val := pathVars[name]; val != "" {
			args = append(args, reflect.ValueOf(val))

			continue
		}

		args = append(args, reflect.Zero(t))
	}

	return args, nil
}
