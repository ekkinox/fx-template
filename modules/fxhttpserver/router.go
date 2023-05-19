package fxhttpserver

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type Middleware interface {
	Handle() echo.MiddlewareFunc
}

type Handler interface {
	Handle() echo.HandlerFunc
}

type Router struct {
	handlers              []Handler
	handlerDefinitions    []HandlerDefinition
	middlewares           []Middleware
	middlewareDefinitions []MiddlewareDefinition
}

type FxRouterParam struct {
	fx.In
	Handlers              []Handler              `group:"http-server-handlers"`
	HandlerDefinitions    []HandlerDefinition    `group:"http-server-handler-definitions"`
	Middlewares           []Middleware           `group:"http-server-middlewares"`
	MiddlewareDefinitions []MiddlewareDefinition `group:"http-server-middleware-definitions"`
}

func NewFxRouter(p FxRouterParam) *Router {
	return &Router{
		handlers:              p.Handlers,
		handlerDefinitions:    p.HandlerDefinitions,
		middlewares:           p.Middlewares,
		middlewareDefinitions: p.MiddlewareDefinitions,
	}
}

func (r *Router) ResolveHandlers() ([]ResolvedHandler, error) {

	var resolvedHandlers []ResolvedHandler

	for _, handlerDef := range r.handlerDefinitions {

		var resHandler ResolvedHandler

		var resMiddlewares []echo.MiddlewareFunc
		for _, middlewareDef := range handlerDef.Middlewares() {
			if middlewareDef.Concrete() {
				resMiddlewares = append(
					resMiddlewares,
					middlewareDef.Middleware().(echo.MiddlewareFunc),
				)
			} else {
				registeredMiddleware, err := r.lookupRegisteredMiddleware(middlewareDef.Middleware().(string))
				if err != nil {
					return nil, err
				}

				resMiddlewares = append(resMiddlewares, registeredMiddleware.Handle())
			}
		}

		if handlerDef.Concrete() {
			resHandler = newResolvedHandler(
				handlerDef.Method(),
				handlerDef.Path(),
				handlerDef.Handler().(func(echo.Context) error),
				resMiddlewares...,
			)
		} else {
			registeredHandler, err := r.lookupRegisteredHandler(handlerDef.Handler().(string))
			if err != nil {
				return nil, err
			}

			resHandler = newResolvedHandler(
				handlerDef.Method(),
				handlerDef.Path(),
				registeredHandler.Handle(),
				resMiddlewares...,
			)
		}

		resolvedHandlers = append(resolvedHandlers, resHandler)
	}

	return resolvedHandlers, nil
}

func (r *Router) lookupRegisteredMiddleware(middleware string) (Middleware, error) {
	for _, m := range r.middlewares {
		if getType(m) == middleware {
			return m, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("cannot find middleware for type %s", middleware))
}

func (r *Router) lookupRegisteredHandler(handler string) (Handler, error) {
	for _, h := range r.handlers {
		if getType(h) == handler {
			return h, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("cannot find handler for type %s", handler))
}
