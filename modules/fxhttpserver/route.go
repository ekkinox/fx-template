package fxhttpserver

import (
	"errors"
	"fmt"
)

type Route interface {
	Handler() string
	Method() string
	Path() string
}

type route struct {
	handler string
	method  string
	path    string
}

func newRoute(handler string, method string, path string) *route {
	return &route{
		handler: handler,
		method:  method,
		path:    path,
	}
}

func (r *route) Handler() string {
	return r.handler
}

func (r *route) Method() string {
	return r.method
}

func (r *route) Path() string {
	return r.path
}

func getRouteForHandler(routes []Route, handler string) (Route, error) {
	for _, r := range routes {
		if r.Handler() == handler {
			return r, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("cannot find route for handler %s", handler))
}
