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
	middlewares              []Middleware
	middlewareDefinitions    []MiddlewareDefinition
	handlers                 []Handler
	handlerDefinitions       []HandlerDefinition
	handlersGroupDefinitions []HandlersGroupDefinition
}

type FxRouterParam struct {
	fx.In
	Middlewares              []Middleware              `group:"http-server-middlewares"`
	MiddlewareDefinitions    []MiddlewareDefinition    `group:"http-server-middleware-definitions"`
	Handlers                 []Handler                 `group:"http-server-handlers"`
	HandlerDefinitions       []HandlerDefinition       `group:"http-server-handler-definitions"`
	HandlersGroupDefinitions []HandlersGroupDefinition `group:"http-server-handlers-group-definitions"`
}

func NewFxRouter(p FxRouterParam) *Router {
	return &Router{
		middlewares:              p.Middlewares,
		middlewareDefinitions:    p.MiddlewareDefinitions,
		handlers:                 p.Handlers,
		handlerDefinitions:       p.HandlerDefinitions,
		handlersGroupDefinitions: p.HandlersGroupDefinitions,
	}
}

func (r *Router) ResolveMiddlewares() ([]ResolvedMiddleware, error) {
	var resolvedMiddlewares []ResolvedMiddleware

	for _, middlewareDef := range r.middlewareDefinitions {

		var resolvedMiddleware ResolvedMiddleware

		if middlewareDef.Kind() != Attached {
			if middlewareDef.Concrete() {
				resolvedMiddleware = newResolvedMiddleware(
					middlewareDef.Middleware().(func(echo.HandlerFunc) echo.HandlerFunc),
					middlewareDef.Kind(),
				)
			} else {
				registeredMiddleware, err := r.lookupRegisteredMiddleware(middlewareDef.Middleware().(string))
				if err != nil {
					return nil, err
				}

				resolvedMiddleware = newResolvedMiddleware(
					registeredMiddleware.Handle(),
					middlewareDef.Kind(),
				)
			}
		}

		resolvedMiddlewares = append(resolvedMiddlewares, resolvedMiddleware)
	}

	return resolvedMiddlewares, nil
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

func (r *Router) ResolveHandlersGroups() ([]ResolvedHandlersGroup, error) {

	var resolvedHandlersGroups []ResolvedHandlersGroup

	for _, handlerGroupDef := range r.handlersGroupDefinitions {

		var groupResolvedMiddlewares []echo.MiddlewareFunc
		for _, middlewareDef := range handlerGroupDef.Middlewares() {
			if middlewareDef.Concrete() {
				groupResolvedMiddlewares = append(
					groupResolvedMiddlewares,
					middlewareDef.Middleware().(echo.MiddlewareFunc),
				)
			} else {
				registeredMiddleware, err := r.lookupRegisteredMiddleware(middlewareDef.Middleware().(string))
				if err != nil {
					return nil, err
				}

				groupResolvedMiddlewares = append(groupResolvedMiddlewares, registeredMiddleware.Handle())
			}
		}

		var groupResolvedHandlers []ResolvedHandler
		for _, handlerDef := range handlerGroupDef.Handlers() {

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

			groupResolvedHandlers = append(groupResolvedHandlers, resHandler)
		}

		resolvedHandlersGroups = append(
			resolvedHandlersGroups,
			newResolvedHandlersGroup(
				handlerGroupDef.Prefix(),
				groupResolvedHandlers,
				groupResolvedMiddlewares...,
			),
		)
	}

	return resolvedHandlersGroups, nil
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
