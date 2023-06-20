package fxhttpserver

import (
	"github.com/ekkinox/fx-template/modules/fxauthenticationcontext"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

func applyDefaultMiddlewares(e *echo.Echo, c *fxconfig.Config, l *EchoLogger) *echo.Echo {
	// failures recovery
	e.Use(middleware.Recover())

	// x-request-id
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
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

	// open telemetry
	e.Use(otelecho.Middleware(c.AppName()))

	// logging
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
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
			evt := l.logger.Info()
			if v.Error != nil {
				evt = l.logger.Error().Err(v.Error)
			}

			evt.
				Str("method", v.Method).
				Str("uri", v.URI).
				Int("status", v.Status).
				Str("latency", v.Latency.String()).
				Str("x-request-id", v.RequestID).
				Str("traceparent", c.Request().Header.Get(headerTraceParent)).
				Str("remote-ip", v.RemoteIP).
				Str("referer", v.Referer).
				Msg("call")

			return nil
		},
	}))

	// authentication context
	e.Use(fxauthenticationcontext.Middleware())

	return e
}
