package fxhttpserver

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/ekkinox/fx-template/modules/fxauthenticationcontext"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.uber.org/fx"
)

const (
	DefaultPort       = 8080
	HeaderXRequestId  = "x-request-id"
	HeaderTraceParent = "traceparent"
)

var FxHttpServerModule = fx.Module(
	"http-server",
	fx.Provide(
		NewDefaultHttpServerFactory,
		NewFxHttpServerRegistry,
		NewFxHttpServer,
	),
	fx.Invoke(func(*echo.Echo) {}),
)

type FxHttpServerParam struct {
	fx.In
	LifeCycle     fx.Lifecycle
	Factory       HttpServerFactory
	Registry      *HttpServerRegistry
	Config        *fxconfig.Config
	Logger        *fxlogger.Logger
	HealthChecker *fxhealthchecker.HealthChecker
}

func NewFxHttpServer(p FxHttpServerParam) (*echo.Echo, error) {
	// logger
	echoLogger := NewEchoLogger(p.Logger)

	// server
	httpServer, err := p.Factory.Create(
		WithDebug(p.Config.AppDebug()),
		WithBanner(false),
		WithLogger(echoLogger),
		WithHttpErrorHandler(NewHttpServerErrorHandler(p.Config)),
	)
	if err != nil {
		p.Logger.Error().Err(err).Msg("failed to create http server")

		return nil, err
	}

	// default middlewares
	httpServer = withDefaultMiddlewares(httpServer, p.Config, echoLogger)

	// registry resources
	httpServer = withRegisteredResources(httpServer, p.Registry, echoLogger)

	// debuggers
	httpServer = withDebuggers(httpServer, p.Config, echoLogger)

	// health checker
	httpServer = withHealthChecker(httpServer, p.HealthChecker)

	// lifecycles
	p.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			port := p.Config.GetInt("http.server.port")
			if port == 0 {
				port = DefaultPort
			}

			go httpServer.Start(fmt.Sprintf(":%d", port))

			return nil

		},
		OnStop: func(ctx context.Context) error {
			return httpServer.Shutdown(ctx)
		},
	})

	return httpServer, nil
}

func withDefaultMiddlewares(httpServer *echo.Echo, config *fxconfig.Config, echoLogger *EchoLogger) *echo.Echo {
	// recovery
	httpServer.Use(middleware.Recover())

	// open telemetry
	httpServer.Use(otelecho.Middleware(config.AppName()))

	// request-id
	httpServer.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = generateRequestId()
				req.Header.Set(echo.HeaderXRequestID, rid)
			}
			res.Header().Set(echo.HeaderXRequestID, rid)

			return next(c)
		}
	})

	// logger
	httpServer.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:    true,
		LogURI:       true,
		LogStatus:    true,
		LogRequestID: true,
		LogLatency:   true,
		LogUserAgent: true,
		LogRemoteIP:  true,
		LogReferer:   true,
		LogError:     true,
		BeforeNextFunc: func(c echo.Context) {
			requestId := c.Request().Header.Get(HeaderXRequestId)
			if requestId == "" {
				requestId = c.Response().Header().Get(HeaderXRequestId)
			}

			traceParent := c.Request().Header.Get(HeaderTraceParent)
			if traceParent == "" {
				traceParent = c.Response().Header().Get(HeaderTraceParent)
			}

			newLogger := echoLogger.
				logger.
				With().
				Str(HeaderXRequestId, requestId).
				Str(HeaderTraceParent, traceParent).
				Logger()

			c.SetRequest(c.Request().WithContext(newLogger.WithContext(c.Request().Context())))
			c.SetLogger(NewEchoLogger(fxlogger.FromZerolog(newLogger)))
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			evt := echoLogger.logger.Info()
			if v.Error != nil {
				evt = echoLogger.logger.Error().Err(v.Error)
			}

			evt.
				Str("method", v.Method).
				Str("uri", v.URI).
				Int("status", v.Status).
				Str("latency", v.Latency.String()).
				Str("x-request-id", v.RequestID).
				Str("traceparent", c.Request().Header.Get(HeaderTraceParent)).
				Str("remote-ip", v.RemoteIP).
				Str("referer", v.Referer).
				Msg("call")

			return nil
		},
	}))

	// auth context
	httpServer.Use(fxauthenticationcontext.Middleware())

	return httpServer
}

func withRegisteredResources(httpServer *echo.Echo, registry *HttpServerRegistry, echoLogger *EchoLogger) *echo.Echo {
	// groups
	resolvedHandlersGroups, err := registry.ResolveHandlersGroups()
	if err != nil {
		echoLogger.Error("cannot resolve router handlers groups: %v", err)
	}

	for _, g := range resolvedHandlersGroups {
		group := httpServer.Group(g.Prefix(), g.Middlewares()...)
		for _, h := range g.Handlers() {
			group.Add(
				strings.ToUpper(h.Method()),
				h.Path(),
				h.Handler(),
				h.Middlewares()...,
			)
			echoLogger.Debugf("registering handler in group for [%s]%s%s", h.Method(), g.Prefix(), h.Path())
		}
		echoLogger.Debugf("registered handlers group for prefix %s", g.Prefix())
	}

	// middlewares
	resolvedMiddlewares, err := registry.ResolveMiddlewares()
	if err != nil {
		echoLogger.Error("cannot resolve router middlewares: %v", err)
	}

	for _, m := range resolvedMiddlewares {
		if m.Kind() == GlobalPre {
			httpServer.Pre(m.Middleware())
		}
		if m.Kind() == GlobalUse {
			httpServer.Use(m.Middleware())
		}
		echoLogger.Debugf("registered %s middleware %T", m.Kind().String(), m.Middleware())
	}

	// handlers
	resolvedHandlers, err := registry.ResolveHandlers()
	if err != nil {
		echoLogger.Error("cannot resolve router handlers: %v", err)
	}

	for _, h := range resolvedHandlers {
		httpServer.Add(
			strings.ToUpper(h.Method()),
			h.Path(),
			h.Handler(),
			h.Middlewares()...,
		)
		echoLogger.Debugf("registered handler for [%s]%s", h.Method(), h.Path())
	}

	return httpServer
}

func withDebuggers(httpServer *echo.Echo, config *fxconfig.Config, echoLogger *EchoLogger) *echo.Echo {
	if config.AppDebug() {

		echoLogger.Debugf("enabling debugging endpoints")

		debugGroup := httpServer.Group("/_debug")

		// configs
		debugGroup.GET("/config", func(c echo.Context) error {
			return c.JSON(http.StatusOK, config.AllSettings())
		})

		// routes
		debugGroup.GET("/routes", func(c echo.Context) error {
			return c.JSON(http.StatusOK, httpServer.Routes())
		})

		// version
		debugGroup.GET("/version", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{
				"application": config.AppName(),
				"version":     config.AppVersion(),
			})
		})
	}

	return httpServer
}

func withHealthChecker(httpServer *echo.Echo, healthChecker *fxhealthchecker.HealthChecker) *echo.Echo {
	httpServer.GET("/_health", func(c echo.Context) error {
		r := healthChecker.Run(c.Request().Context())

		status := http.StatusOK
		if !r.Success {
			status = http.StatusInternalServerError
		}

		return c.JSON(status, r)
	})

	return httpServer
}
