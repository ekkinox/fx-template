package fxhttpserver

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"reflect"
)

type Handler interface {
	Handle() echo.HandlerFunc
}

func RegisterHandler(method string, path string, constructor any) fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(
				constructor,
				fx.As(new(Handler)),
				fx.ResultTags(`group:"http-server-handlers"`),
			),
		),
		fx.Supply(
			fx.Annotate(
				newRoute(reflect.TypeOf(constructor).Out(0).String(), method, path),
				fx.As(new(Route)),
				fx.ResultTags(`group:"http-server-routes"`),
			),
		),
	)
}
