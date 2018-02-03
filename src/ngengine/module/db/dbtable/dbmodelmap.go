package dbtable

import (
	"errors"
	"reflect"
)

var DbtableMap map[string]reflect.Type

func init() {
	DbtableMap = make(map[string]reflect.Type)
}

func DbstructPool(structName string) (interface{}, error) {
	if j, ok := DbtableMap[structName]; ok {
		return reflect.New(j).Interface(), nil
	}

	return nil, errors.New("dont have this table")
}
