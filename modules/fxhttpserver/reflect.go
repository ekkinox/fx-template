package fxhttpserver

import (
	"github.com/labstack/echo/v4"
	"reflect"
)

func getType(target any) string {
	return reflect.TypeOf(target).String()
}

func getReturnType(target any) string {
	return reflect.TypeOf(target).Out(0).String()
}

func isConcreteMiddleware(middleware any) bool {
	return reflect.TypeOf(middleware).ConvertibleTo(reflect.TypeOf(echo.MiddlewareFunc(nil)))
}

func isConcreteHandler(handler any) bool {
	return reflect.TypeOf(handler).ConvertibleTo(reflect.TypeOf(echo.HandlerFunc(nil)))
}
