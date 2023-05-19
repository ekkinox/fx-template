package fxhttpserver

import "github.com/labstack/echo/v4"

type ResolvedHandler interface {
	Method() string
	Path() string
	Handler() echo.HandlerFunc
	Middlewares() []echo.MiddlewareFunc
}

type resolvedHandler struct {
	method      string
	path        string
	handler     echo.HandlerFunc
	middlewares []echo.MiddlewareFunc
}

func newResolvedHandler(method string, path string, handler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) *resolvedHandler {
	return &resolvedHandler{
		method:      method,
		path:        path,
		handler:     handler,
		middlewares: middlewares,
	}
}

func (r *resolvedHandler) Method() string {
	return r.method
}

func (r *resolvedHandler) Path() string {
	return r.path
}

func (r *resolvedHandler) Handler() echo.HandlerFunc {
	return r.handler
}

func (r *resolvedHandler) Middlewares() []echo.MiddlewareFunc {
	return r.middlewares
}
