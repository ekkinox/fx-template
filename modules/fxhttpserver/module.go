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
	"go.opentelemetry.io/otel/sdk/trace"
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
)

type FxHttpServerParam struct {
	fx.In
	LifeCycle      fx.Lifecycle
	Factory        HttpServerFactory
	Registry       *HttpServerRegistry
	Config         *fxconfig.Config
	Logger         *fxlogger.Logger
	TracerProvider *trace.TracerProvider
	HealthChecker  *fxhealthchecker.HealthChecker
}

func StartFxHttpServer() fx.Option {
	return fx.Invoke(func(*echo.Echo) {})
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

	// middlewares
	httpServer = withDefaultMiddlewares(httpServer, p.Config, p.TracerProvider)

	// registry
	httpServer = withRegisteredResources(httpServer, p.Registry)

	// debuggers
	httpServer = withDebugEndpoints(httpServer, p.Config)

	// health checker
	httpServer = withHealthCheckEndpoint(httpServer, p.HealthChecker)

	// lifecycles
	p.LifeCycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			port := p.Config.GetInt("modules.http.server.port")
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

func withDefaultMiddlewares(httpServer *echo.Echo, config *fxconfig.Config, tracerProvider *trace.TracerProvider) *echo.Echo {

	// recovery
	httpServer.Use(middleware.Recover())

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
		LogLatency:   true,
		LogUserAgent: true,
		LogRemoteIP:  true,
		LogReferer:   true,
		LogError:     true,
		LogHeaders: []string{
			HeaderXRequestId,
			HeaderTraceParent,
		},
		BeforeNextFunc: func(c echo.Context) {

			fields := map[string]interface{}{}

			for _, h := range []string{
				HeaderXRequestId,
				HeaderTraceParent,
			} {
				v := extractHeader(c, h)
				if v != nil {
					fields[h] = v
				}
			}

			cl := c.Logger().(*EchoLogger).ToZerolog().With().Fields(fields).Logger()

			c.SetRequest(c.Request().WithContext(cl.WithContext(c.Request().Context())))
			c.SetLogger(NewEchoLogger(fxlogger.FromZerolog(cl)))
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {

			zl := c.Logger().(*EchoLogger).ToZerolog()

			evt := zl.Info()
			if v.Error != nil {
				evt = zl.Error().Err(v.Error)
			}

			evt.
				Str("method", v.Method).
				Str("uri", v.URI).
				Int("status", v.Status).
				Str("latency", v.Latency.String()).
				Str("remote-ip", v.RemoteIP).
				Str("referer", v.Referer)

			for headerName, headerValues := range v.Headers {
				value := ""
				if len(headerValues) != 0 {
					value = headerValues[0]
				}

				evt.Str(strings.ToLower(headerName), value)
			}

			evt.Msg("request")

			return nil
		},
	}))

	// open telemetry
	if config.GetBool("modules.http.tracer.enabled") {
		httpServer.Use(otelecho.Middleware(
			config.AppName(),
			otelecho.WithTracerProvider(tracerProvider),
		))
	}

	// auth context
	if config.GetBool("modules.auth.enabled") {
		httpServer.Use(fxauthenticationcontext.Middleware(config.GetBool("modules.auth.blocking")))
	}

	return httpServer
}

func withRegisteredResources(httpServer *echo.Echo, registry *HttpServerRegistry) *echo.Echo {

	// groups
	resolvedHandlersGroups, err := registry.ResolveHandlersGroups()
	if err != nil {
		httpServer.Logger.Errorf("cannot resolve router handlers groups: %v", err)
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
			httpServer.Logger.Debugf("registering handler in group for [%s]%s%s", h.Method(), g.Prefix(), h.Path())
		}
		httpServer.Logger.Debugf("registered handlers group for prefix %s", g.Prefix())
	}

	// middlewares
	resolvedMiddlewares, err := registry.ResolveMiddlewares()
	if err != nil {
		httpServer.Logger.Errorf("cannot resolve router middlewares: %v", err)
	}

	for _, m := range resolvedMiddlewares {
		if m.Kind() == GlobalPre {
			httpServer.Pre(m.Middleware())
		}
		if m.Kind() == GlobalUse {
			httpServer.Use(m.Middleware())
		}
		httpServer.Logger.Debugf("registered %s middleware %T", m.Kind().String(), m.Middleware())
	}

	// handlers
	resolvedHandlers, err := registry.ResolveHandlers()
	if err != nil {
		httpServer.Logger.Errorf("cannot resolve router handlers: %v", err)
	}

	for _, h := range resolvedHandlers {
		httpServer.Add(
			strings.ToUpper(h.Method()),
			h.Path(),
			h.Handler(),
			h.Middlewares()...,
		)
		httpServer.Logger.Debugf("registered handler for [%s]%s", h.Method(), h.Path())
	}

	return httpServer
}

func withDebugEndpoints(httpServer *echo.Echo, config *fxconfig.Config) *echo.Echo {
	if config.AppDebug() {
		httpServer.Logger.Debugf("enabling debugging endpoints")

		debugGroup := httpServer.Group("/_debug")

		// configs endpoint
		debugGroup.GET("/config", func(c echo.Context) error {
			return c.JSON(http.StatusOK, config.AllSettings())
		})

		// routes endpoint
		debugGroup.GET("/routes", func(c echo.Context) error {
			return c.JSON(http.StatusOK, httpServer.Routes())
		})

		// version endpoint
		debugGroup.GET("/version", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{
				"application": config.AppName(),
				"version":     config.AppVersion(),
			})
		})
	}

	return httpServer
}

func withHealthCheckEndpoint(httpServer *echo.Echo, healthChecker *fxhealthchecker.HealthChecker) *echo.Echo {

	httpServer.Logger.Debugf("enabling health check endpoint")

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

func extractHeader(c echo.Context, headerName string) *string {
	headerValue := c.Request().Header.Get(headerName)

	if headerValue == "" {
		headerValue = c.Response().Header().Get(headerName)
	}

	if headerValue == "" {
		return nil
	}

	return &headerValue
}
