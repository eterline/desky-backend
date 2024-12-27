package utils

import (
	"reflect"
)

func ContainsInStruct(obj, value any) bool {
	values := reflect.ValueOf(obj)

	if values.Kind() != reflect.Struct {
		panic("obj must be a struct")
	}

	for i := 0; i < values.NumField(); i++ {
		field := values.Field(i)
		if reflect.DeepEqual(field.Interface(), value) {
			return true
		}
	}

	return false
}

func ContainsInListOfStruct[Type any](obj []Type, value any) bool {
	for _, i := range obj {
		if ContainsInStruct(i, value) {
			return true
		}
	}
	return false
}
