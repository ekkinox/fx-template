package fxhttpserver

import (
	"errors"
	"fmt"
)

type Route interface {
	Method() string
	Path() string
	Handler() string
	Middlewares() []string
}

type route struct {
	method      string
	path        string
	handler     string
	middlewares []string
}

func newRoute(method string, path string, handler string, middlewares ...string) *route {
	return &route{
		method:      method,
		path:        path,
		handler:     handler,
		middlewares: middlewares,
	}
}

func (r *route) Method() string {
	return r.method
}

func (r *route) Path() string {
	return r.path
}

func (r *route) Handler() string {
	return r.handler
}

func (r *route) Middlewares() []string {
	return r.middlewares
}

func findRouteForHandler(routes []Route, handler string) (Route, error) {
	for _, r := range routes {
		if r.Handler() == handler {
			return r, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("cannot find route for handler %s", handler))
}
