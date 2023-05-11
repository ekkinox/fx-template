package fxhttpserver

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type HttpServerHandler interface {
	Method() string
	Path() string
	Handler() echo.HandlerFunc
	Middlewares() []echo.MiddlewareFunc
}

func AsHttpServerHandler(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(HttpServerHandler)),
		fx.ResultTags(`group:"http-server-handlers"`),
	)
}
