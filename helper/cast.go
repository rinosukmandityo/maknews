package helper

import (
	"reflect"
)

func Value(o interface{}) reflect.Value {
	return reflect.Indirect(reflect.ValueOf(o))
}

func Kind(o interface{}) reflect.Kind {
	return Value(o).Kind()
}
