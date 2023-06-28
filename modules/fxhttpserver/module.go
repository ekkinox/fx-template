package fxhttpserver

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

const DefaultPort = 8080

var FxHttpServerModule = fx.Module(
	"http-server",
	// modules dependencies
	fxhealthchecker.FxHealthCheckerModule,
	// http server
	fx.Provide(
		NewFxRouter,
		NewFxHttpServer,
	),
	fx.Invoke(func(*echo.Echo) {}),
)

type FxHttpServerParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Config    *fxconfig.Config
	Logger    *fxlogger.Logger
	Checker   *fxhealthchecker.Checker
	Router    *Router
}

func NewFxHttpServer(p FxHttpServerParam) *echo.Echo {
	// logger
	l := NewEchoLogger(p.Logger)

	// echo
	e := echo.New()
	e.HideBanner = true
	e.Debug = p.Config.AppDebug()
	e.Logger = l
	e.HTTPErrorHandler = NewHttpServerErrorHandler(p.Config)

	// middlewares
	e = applyDefaultMiddlewares(e, p.Config, l)

	// groups
	resolvedHandlersGroups, err := p.Router.ResolveHandlersGroups()
	if err != nil {
		l.Error("cannot resolve router handlers groups: %v", err)
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
			l.Debugf("registering handler in group for [%s]%s%s", h.Method(), g.Prefix(), h.Path())
		}
		l.Debugf("registered handlers group for prefix %s", g.Prefix())
	}

	// middlewares
	resolvedMiddlewares, err := p.Router.ResolveMiddlewares()
	if err != nil {
		l.Error("cannot resolve router middlewares: %v", err)
	}

	for _, m := range resolvedMiddlewares {
		if m.Kind() == GlobalPre {
			e.Pre(m.Middleware())
		}
		if m.Kind() == GlobalUse {
			e.Use(m.Middleware())
		}
		l.Debugf("registered %s middleware %T", m.Kind().String(), m.Middleware())
	}

	// handlers
	resolvedHandlers, err := p.Router.ResolveHandlers()
	if err != nil {
		l.Error("cannot resolve router handlers: %v", err)
	}

	for _, h := range resolvedHandlers {
		e.Add(
			strings.ToUpper(h.Method()),
			h.Path(),
			h.Handler(),
			h.Middlewares()...,
		)
		l.Debugf("registered handler for [%s]%s", h.Method(), h.Path())
	}

	// debuggers
	if p.Config.AppDebug() {
		g := e.Group("/_debug")
		// configs
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

	// health checker
	e.GET("/_health", func(c echo.Context) error {
		r := p.Checker.Run(c.Request().Context())

		status := http.StatusOK
		if !r.Success {
			status = http.StatusInternalServerError
		}

		return c.JSON(status, r)
	})

	// lifecycles
	p.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			port := p.Config.GetInt("http.server.port")
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
