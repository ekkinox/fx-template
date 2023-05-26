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
	// logger
	l := NewEchoLogger(p.Logger)

	// echo
	e := echo.New()
	e.HideBanner = true
	e.Debug = p.Config.AppDebug()
	e.Logger = l

	// middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(otelecho.Middleware(p.Config.AppName()))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:    true,
		LogURI:       true,
		LogStatus:    true,
		LogRequestID: true,
		LogLatency:   true,
		LogUserAgent: true,
		LogRemoteIP:  true,
		LogReferer:   true,
		BeforeNextFunc: func(c echo.Context) {
			requestId := c.Request().Header.Get(headerRequestID)
			if requestId == "" {
				requestId = c.Response().Header().Get(headerRequestID)
			}

			traceParent := c.Request().Header.Get(headerTraceParent)
			if traceParent == "" {
				traceParent = c.Response().Header().Get(headerTraceParent)
			}

			corrLogger := l.logger.
				With().
				Str(headerRequestID, requestId).
				Str(headerTraceParent, traceParent).
				Logger()

			c.SetRequest(c.Request().WithContext(corrLogger.WithContext(c.Request().Context())))
			c.SetLogger(NewEchoLogger(fxlogger.FromLogger(corrLogger)))
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			l.logger.Info().
				Str("method", v.Method).
				Str("uri", v.URI).
				Int("status", v.Status).
				Str("latency", v.Latency.String()).
				Str("x-request-id", v.RequestID).
				Str("traceparent", c.Request().Header.Get(headerTraceParent)).
				Str("remote-ip", v.RemoteIP).
				Str("referer", v.Referer).
				Msg("")

			return nil
		},
	}))

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

	// health checker
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
