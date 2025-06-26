package core

import (
	"encoding/json"
	"net/http"
	"reflect"
)

func BindJSONBody(req *http.Request, targetType reflect.Type) (any, error) {
	ptr := reflect.New(targetType).Interface()
	err := json.NewDecoder(req.Body).Decode(ptr)
	if err != nil {
		return nil, err
	}
	return reflect.Indirect(reflect.ValueOf(ptr)).Interface(), nil
}
