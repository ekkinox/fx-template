package fxhttpserver

import (
	"go.uber.org/fx"
)

func RegisterHandler(method string, path string, handler any, middlewares ...any) fx.Option {

	var providers []any

	var middlewareDefs []MiddlewareDefinition
	for _, middleware := range middlewares {
		if !isConcreteMiddleware(middleware) {
			providers = append(
				providers,
				fx.Annotate(
					middleware,
					fx.As(new(Middleware)),
					fx.ResultTags(`group:"http-server-middlewares"`),
				),
			)

			middlewareDefs = append(middlewareDefs, newMiddlewareDefinition(getReturnType(middleware)))
		} else {
			middlewareDefs = append(middlewareDefs, newMiddlewareDefinition(middleware))
		}
	}

	var handlerDef HandlerDefinition
	if !isConcreteHandler(handler) {
		providers = append(
			providers,
			fx.Annotate(
				handler,
				fx.As(new(Handler)),
				fx.ResultTags(`group:"http-server-handlers"`),
			),
		)
		handlerDef = newHandlerDefinition(method, path, getReturnType(handler), middlewareDefs)
	} else {
		handlerDef = newHandlerDefinition(method, path, handler, middlewareDefs)
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
