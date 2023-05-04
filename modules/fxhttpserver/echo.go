package fxhttpserver

import (
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"strings"
)

type EchoParam struct {
	fx.In
	Logger   *fxlogger.Logger
	Handlers []HttpServerHandler `group:"http-server-handlers"`
}

type EchoResult struct {
	fx.Out
	Echo *echo.Echo
}

func NewEcho(p EchoParam) EchoResult {
	e := echo.New()

	for _, h := range p.Handlers {
		e.Add(strings.ToUpper(h.Method()), h.Path(), h.Handle())
	}

	e.Use(middleware.RequestLoggerWithConfig(
		middleware.RequestLoggerConfig{
			LogURI:    true,
			LogStatus: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				p.Logger.Info().
					Str("URI", v.URI).
					Int("status", v.Status).
					Msg("request")

				return nil
			},
		},
	))

	return EchoResult{
		Echo: e,
	}
}
