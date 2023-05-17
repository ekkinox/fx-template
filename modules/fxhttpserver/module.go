package fxhttpserver

import (
	"context"
	"fmt"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/ekkinox/fx-template/modules/fxtracer"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.uber.org/fx"
	"net/http"
	"reflect"
	"strings"
)

var FxHttpServerModule = fx.Module("http-server",
	// modules dependencies
	fxconfig.FxConfigModule,
	fxlogger.FxLoggerModule,
	fxtracer.FxTracerModule,
	// http server
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
	Handlers  []Handler `group:"http-server-handlers"`
	Routes    []Route   `group:"http-server-routes"`
}

func NewFxHttpServer(p FxHttpServerParam) *echo.Echo {
	// echo
	e := echo.New()
	e.HideBanner = true
	e.Debug = p.Config.AppConfig.Debug
	e.Logger = p.Logger

	// middlewares
	e.Use(middleware.RequestID())
	e.Use(otelecho.Middleware(p.Config.AppConfig.Name))
	e.Use(fxlogger.Middleware(fxlogger.Config{
		Logger:      p.Logger,
		HandleError: true,
	}))

	// handlers groups
	/*	for _, hg := range p.HandlersGroups {
		g := e.Group(hg.Prefix(), hg.Middlewares()...)
		for _, h := range hg.Handlers() {
			g.Add(strings.ToUpper(h.Method()), h.Path(), h.Handler(), h.Middlewares()...)
		}
	}*/

	// handlers
	for _, h := range p.Handlers {

		r, err := getRouteForHandler(p.Routes, reflect.TypeOf(h).String())
		if err != nil {
			p.Logger.Error("cannot register handler")
		}

		e.Add(strings.ToUpper(r.Method()), r.Path(), h.Handle())
		p.Logger.Infof("registered handler for %s %s", r.Method(), r.Path())
	}

	// debugger
	if p.Config.AppConfig.Debug {
		g := e.Group("/_debug")
		// routes
		g.GET("/routes", func(c echo.Context) error {
			return c.JSON(http.StatusOK, e.Routes())
		})
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
