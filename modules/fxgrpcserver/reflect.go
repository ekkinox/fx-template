package fxgrpcserver

import (
	"reflect"
)

func getType(target any) string {
	return reflect.TypeOf(target).String()
}

func getReturnType(target any) string {
	return reflect.TypeOf(target).Out(0).String()
}
