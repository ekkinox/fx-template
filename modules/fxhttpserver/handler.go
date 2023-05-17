package fxhttpserver

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"reflect"
)

type Middleware interface {
	Handle() echo.MiddlewareFunc
}

type Handler interface {
	Handle() echo.HandlerFunc
}

func RegisterHandler(method string, path string, handler any, middlewares ...any) fx.Option {

	hp := fx.Provide(
		fx.Annotate(
			handler,
			fx.As(new(Handler)),
			fx.ResultTags(`group:"http-server-handlers"`),
		),
	)

	var mm []any
	var mt []string
	var mi []echo.MiddlewareFunc

	for _, m := range middlewares {
		if reflect.TypeOf(m).String() == "echo.MiddlewareFunc" {
			mi = append(
				mi,
				m.(echo.MiddlewareFunc),
			)
		} else {
			mm = append(
				mm,
				fx.Annotate(
					m,
					fx.As(new(Middleware)),
					fx.ResultTags(`group:"http-server-middlewares"`),
				),
			)
			mt = append(
				mt,
				reflect.TypeOf(m).Out(0).String(),
			)
		}

	}
	mp := fx.Provide(mm...)

	rp := fx.Supply(
		fx.Annotate(
			newRoute(method, path, reflect.TypeOf(handler).Out(0).String(), mt, mi),
			fx.As(new(Route)),
			fx.ResultTags(`group:"http-server-routes"`),
		),
	)

	return fx.Options(
		hp,
		mp,
		rp,
	)
}
