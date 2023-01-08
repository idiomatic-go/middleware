package gotp

import "reflect"

func IsNillable(a any) bool {
	return IsPointer(a) || IsPointerType(a)
}

func IsPointer(a any) bool {
	if a == nil {
		return false
	}
	if reflect.TypeOf(a).Kind() != reflect.Pointer {
		return false
	}
	return true
}

func IsPointerType(a any) bool {
	if a == nil {
		return false
	}
	switch reflect.ValueOf(a).Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map:
		return true
	}
	return false
}

// IsNil - determine if the interface{} holds is nil
func IsNil(a any) bool {
	if a == nil {
		return true
	}
	if !IsNillable(a) {
		return false
	}
	return reflect.ValueOf(a).IsNil()
}
