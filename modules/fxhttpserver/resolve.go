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

type ResolvedHandlersGroup interface {
	Prefix() string
	Handlers() []ResolvedHandler
	Middlewares() []echo.MiddlewareFunc
}

type resolvedHandlersGroup struct {
	prefix      string
	handlers    []ResolvedHandler
	middlewares []echo.MiddlewareFunc
}

func newResolvedHandlersGroup(prefix string, handlers []ResolvedHandler, middlewares ...echo.MiddlewareFunc) *resolvedHandlersGroup {
	return &resolvedHandlersGroup{
		prefix:      prefix,
		handlers:    handlers,
		middlewares: middlewares,
	}
}

func (r *resolvedHandlersGroup) Prefix() string {
	return r.prefix
}

func (r *resolvedHandlersGroup) Handlers() []ResolvedHandler {
	return r.handlers
}

func (r *resolvedHandlersGroup) Middlewares() []echo.MiddlewareFunc {
	return r.middlewares
}
