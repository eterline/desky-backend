package utils

import (
	"encoding/json"
	"reflect"
)

func RemoveSliceIndex[Type any](s []Type, index int) []Type {
	return append(s[:index], s[index+1:]...)
}

func PrettyString(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

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

func CutList[Type any](v []Type, min, max int) []Type {
	return v[min:max]
}
