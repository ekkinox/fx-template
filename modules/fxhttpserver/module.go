package fxhttpserver

import (
	"context"
	"fmt"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.uber.org/fx"
	"strings"
)

var FxHttpServerModule = fx.Module("http-server",
	fx.Provide(
		NewFxHttpServer,
	),
	fx.Invoke(func(*echo.Echo) {}),
)

type FxHttpServerParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Config    *fxconfig.Config
	Logger    *fxlogger.Logger
	Handlers  []HttpServerHandler `group:"http-server-handlers"`
}

func NewFxHttpServer(p FxHttpServerParam) *echo.Echo {
	// echo init
	e := echo.New()
	e.HideBanner = true
	e.Logger = p.Logger

	// middlewares
	e.Use(middleware.RequestID())
	e.Use(otelecho.Middleware(p.Config.AppConfig.Name))
	e.Use(fxlogger.Middleware(fxlogger.Config{
		Logger:      p.Logger,
		HandleError: true,
	}))

	// handlers
	for _, h := range p.Handlers {
		e.Add(strings.ToUpper(h.Method()), h.Path(), h.Handler(), h.Middlewares()...)
	}

	// lifecycles
	p.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go e.Start(fmt.Sprintf(":%d", p.Config.AppConfig.Port))
			return nil

		},
		OnStop: func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})

	return e
}
