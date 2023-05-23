package fxhttpserver

import (
	"context"
	"fmt"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/ekkinox/fx-template/modules/fxtracer"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.uber.org/fx"
	"net/http"
	"strings"
)

const DefaultPort = 8080

var FxHttpServerModule = fx.Module("http-server",
	// modules dependencies
	fxconfig.FxConfigModule,
	fxlogger.FxLoggerModule,
	fxtracer.FxTracerModule,
	fxhealthchecker.FxHealthCheckerModule,
	// http server
	fx.Provide(
		NewFxRouter,
		NewFxHttpServer,
	),
	fx.Invoke(func(*Router) {}),
	fx.Invoke(func(*echo.Echo) {}),
)

type FxHttpServerParam struct {
	fx.In
	LifeCycle    fx.Lifecycle
	Config       *fxconfig.Config
	Logger       *fxlogger.Logger
	HeathChecker *fxhealthchecker.Checker
	Router       *Router
}

func NewFxHttpServer(p FxHttpServerParam) *echo.Echo {
	// echo
	e := echo.New()
	e.HideBanner = true
	e.Debug = p.Config.AppDebug()
	e.Logger = p.Logger

	// middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(otelecho.Middleware(p.Config.AppName()))
	e.Use(fxlogger.Middleware(fxlogger.Config{
		Logger:      p.Logger,
		HandleError: true,
	}))

	// groups
	resolvedHandlersGroups, err := p.Router.ResolveHandlersGroups()
	if err != nil {
		p.Logger.Error("cannot resolve router handlers groups: %v", err)
	}

	for _, g := range resolvedHandlersGroups {
		group := e.Group(g.Prefix(), g.Middlewares()...)
		for _, h := range g.Handlers() {
			group.Add(
				strings.ToUpper(h.Method()),
				h.Path(),
				h.Handler(),
				h.Middlewares()...,
			)
			p.Logger.Infof("registering handler in group for [%s]%s%s", h.Method(), g.Prefix(), h.Path())
		}
		p.Logger.Infof("registered handlers group for prefix %s", g.Prefix())
	}

	// handlers
	resolvedHandlers, err := p.Router.ResolveHandlers()
	if err != nil {
		p.Logger.Error("cannot resolve router handlers: %v", err)
	}

	for _, h := range resolvedHandlers {

		e.Add(
			strings.ToUpper(h.Method()),
			h.Path(),
			h.Handler(),
			h.Middlewares()...,
		)

		p.Logger.Infof("registered handler for [%s]%s", h.Method(), h.Path())
	}

	// debugger
	if p.Config.AppDebug() {
		g := e.Group("/_debug")
		// config
		g.GET("/config", func(c echo.Context) error {
			return c.JSON(http.StatusOK, p.Config.AllSettings())
		})
		// routes
		g.GET("/routes", func(c echo.Context) error {
			return c.JSON(http.StatusOK, e.Routes())
		})
		// version
		g.GET("/version", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"version": "0.1.0"})
		})
	}

	// health check
	e.GET("/_health", func(c echo.Context) error {
		r := p.HeathChecker.Run(c.Request().Context())

		status := http.StatusOK
		if !r.Success {
			status = http.StatusInternalServerError
		}

		return c.JSON(status, r)
	})

	// lifecycles
	p.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			port := p.Config.GetInt("app.port")
			if port == 0 {
				port = DefaultPort
			}

			go e.Start(fmt.Sprintf(":%d", port))

			return nil

		},
		OnStop: func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})

	return e
}
