package fxhttpserver

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
)

type Route interface {
	Method() string
	Path() string
	HandlerType() string
	MiddlewareTypes() []string
	MiddlewareInstances() []echo.MiddlewareFunc
}

type route struct {
	method              string
	path                string
	handlerType         string
	middlewareTypes     []string
	middlewareInstances []echo.MiddlewareFunc
}

func newRoute(method string, path string, handlerType string, middlewareTypes []string, middlewareInstances []echo.MiddlewareFunc) *route {
	return &route{
		method:              method,
		path:                path,
		handlerType:         handlerType,
		middlewareTypes:     middlewareTypes,
		middlewareInstances: middlewareInstances,
	}
}

func (r *route) Method() string {
	return r.method
}

func (r *route) Path() string {
	return r.path
}

func (r *route) HandlerType() string {
	return r.handlerType
}

func (r *route) MiddlewareTypes() []string {
	return r.middlewareTypes
}

func (r *route) MiddlewareInstances() []echo.MiddlewareFunc {
	return r.middlewareInstances
}

func findRouteForHandlerType(routes []Route, handlerType string) (Route, error) {
	for _, r := range routes {
		if r.HandlerType() == handlerType {
			return r, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("cannot find route for handler type %s", handlerType))
}
