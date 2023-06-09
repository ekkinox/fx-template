package fxhttpserver

import (
	"go.uber.org/fx"
)

type MiddlewareRegistration struct {
	middleware any
	kind       MiddlewareKind
}

func (m *MiddlewareRegistration) Middleware() any {
	return m.middleware
}

func (m *MiddlewareRegistration) Kind() MiddlewareKind {
	return m.kind
}

func NewMiddlewareRegistration(middleware any, kind MiddlewareKind) *MiddlewareRegistration {
	return &MiddlewareRegistration{
		middleware: middleware,
		kind:       kind,
	}
}

func AsMiddleware(middleware any, kind MiddlewareKind) fx.Option {
	return RegisterMiddleware(NewMiddlewareRegistration(middleware, kind))
}

func RegisterMiddleware(middlewareRegistration *MiddlewareRegistration) fx.Option {

	var providers []any

	var middlewareDef MiddlewareDefinition
	if !isConcreteMiddleware(middlewareRegistration.Middleware()) {
		providers = append(
			providers,
			fx.Annotate(
				middlewareRegistration.Middleware(),
				fx.As(new(Middleware)),
				fx.ResultTags(`group:"http-server-middlewares"`),
			),
		)

		middlewareDef = newMiddlewareDefinition(getReturnType(middlewareRegistration.Middleware()), middlewareRegistration.kind)
	} else {
		middlewareDef = newMiddlewareDefinition(middlewareRegistration.Middleware(), middlewareRegistration.kind)
	}

	return fx.Options(
		fx.Provide(providers...),
		fx.Supply(
			fx.Annotate(
				middlewareDef,
				fx.As(new(MiddlewareDefinition)),
				fx.ResultTags(`group:"http-server-middleware-definitions"`),
			),
		),
	)
}

type HandlerRegistration struct {
	method      string
	path        string
	handler     any
	middlewares []any
}

func NewHandlerRegistration(method string, path string, handler any, middlewares ...any) *HandlerRegistration {
	return &HandlerRegistration{
		method:      method,
		path:        path,
		handler:     handler,
		middlewares: middlewares,
	}
}

func (h *HandlerRegistration) Method() string {
	return h.method
}

func (h *HandlerRegistration) Path() string {
	return h.path
}

func (h *HandlerRegistration) Handler() any {
	return h.handler
}

func (h *HandlerRegistration) Middlewares() []any {
	return h.middlewares
}

func AsHandler(method string, path string, handler any, middlewares ...any) fx.Option {
	return RegisterHandler(NewHandlerRegistration(method, path, handler, middlewares...))
}

func RegisterHandler(handlerRegistration *HandlerRegistration) fx.Option {

	var providers []any

	var middlewareDefs []MiddlewareDefinition
	for _, middleware := range handlerRegistration.Middlewares() {
		if !isConcreteMiddleware(middleware) {
			providers = append(
				providers,
				fx.Annotate(
					middleware,
					fx.As(new(Middleware)),
					fx.ResultTags(`group:"http-server-middlewares"`),
				),
			)

			middlewareDefs = append(middlewareDefs, newMiddlewareDefinition(getReturnType(middleware), Attached))
		} else {
			middlewareDefs = append(middlewareDefs, newMiddlewareDefinition(middleware, Attached))
		}
	}

	var handlerDef HandlerDefinition
	if !isConcreteHandler(handlerRegistration.Handler()) {
		providers = append(
			providers,
			fx.Annotate(
				handlerRegistration.Handler(),
				fx.As(new(Handler)),
				fx.ResultTags(`group:"http-server-handlers"`),
			),
		)
		handlerDef = newHandlerDefinition(
			handlerRegistration.Method(),
			handlerRegistration.Path(),
			getReturnType(handlerRegistration.Handler()),
			middlewareDefs,
		)
	} else {
		handlerDef = newHandlerDefinition(
			handlerRegistration.Method(),
			handlerRegistration.Path(),
			handlerRegistration.Handler(),
			middlewareDefs,
		)
	}

	return fx.Options(
		fx.Provide(providers...),
		fx.Supply(
			fx.Annotate(
				handlerDef,
				fx.As(new(HandlerDefinition)),
				fx.ResultTags(`group:"http-server-handler-definitions"`),
			),
		),
	)
}

type HandlersGroupRegistration struct {
	prefix                string
	handlersRegistrations []*HandlerRegistration
	middlewares           []any
}

func NewHandlersGroupRegistration(prefix string, handlersRegistrations []*HandlerRegistration, middlewares ...any) *HandlersGroupRegistration {
	return &HandlersGroupRegistration{
		prefix:                prefix,
		handlersRegistrations: handlersRegistrations,
		middlewares:           middlewares,
	}
}

func (h *HandlersGroupRegistration) Prefix() string {
	return h.prefix
}

func (h *HandlersGroupRegistration) HandlersRegistrations() []*HandlerRegistration {
	return h.handlersRegistrations
}

func (h *HandlersGroupRegistration) Middlewares() []any {
	return h.middlewares
}

func AsHandlersGroup(prefix string, handlersRegistrations []*HandlerRegistration, middlewares ...any) fx.Option {
	return RegisterHandlersGroup(NewHandlersGroupRegistration(prefix, handlersRegistrations, middlewares...))
}

func RegisterHandlersGroup(handlersGroupRegistration *HandlersGroupRegistration) fx.Option {
	var providers []any

	var groupMiddlewareDefs []MiddlewareDefinition
	for _, middleware := range handlersGroupRegistration.Middlewares() {
		if !isConcreteMiddleware(middleware) {
			providers = append(
				providers,
				fx.Annotate(
					middleware,
					fx.As(new(Middleware)),
					fx.ResultTags(`group:"http-server-middlewares"`),
				),
			)

			groupMiddlewareDefs = append(groupMiddlewareDefs, newMiddlewareDefinition(getReturnType(middleware), Attached))
		} else {
			groupMiddlewareDefs = append(groupMiddlewareDefs, newMiddlewareDefinition(middleware, Attached))
		}
	}

	var groupHandlerDefs []HandlerDefinition
	for _, handlerRegistration := range handlersGroupRegistration.HandlersRegistrations() {

		var handlerDef HandlerDefinition
		var middlewareDefs []MiddlewareDefinition

		for _, middleware := range handlerRegistration.Middlewares() {
			if !isConcreteMiddleware(middleware) {
				providers = append(
					providers,
					fx.Annotate(
						middleware,
						fx.As(new(Middleware)),
						fx.ResultTags(`group:"http-server-middlewares"`),
					),
				)

				middlewareDefs = append(middlewareDefs, newMiddlewareDefinition(getReturnType(middleware), Attached))
			} else {
				middlewareDefs = append(middlewareDefs, newMiddlewareDefinition(middleware, Attached))
			}
		}

		if !isConcreteHandler(handlerRegistration.Handler()) {
			providers = append(
				providers,
				fx.Annotate(
					handlerRegistration.Handler(),
					fx.As(new(Handler)),
					fx.ResultTags(`group:"http-server-handlers"`),
				),
			)
			handlerDef = newHandlerDefinition(
				handlerRegistration.Method(),
				handlerRegistration.Path(),
				getReturnType(handlerRegistration.Handler()),
				middlewareDefs,
			)
		} else {
			handlerDef = newHandlerDefinition(
				handlerRegistration.Method(),
				handlerRegistration.Path(),
				handlerRegistration.Handler(),
				middlewareDefs,
			)
		}

		groupHandlerDefs = append(groupHandlerDefs, handlerDef)
	}

	handlersGroupDef := newHandlersGroupDefinition(
		handlersGroupRegistration.Prefix(),
		groupHandlerDefs,
		groupMiddlewareDefs,
	)

	return fx.Options(
		fx.Provide(providers...),
		fx.Supply(
			fx.Annotate(
				handlersGroupDef,
				fx.As(new(HandlersGroupDefinition)),
				fx.ResultTags(`group:"http-server-handlers-group-definitions"`),
			),
		),
	)
}
