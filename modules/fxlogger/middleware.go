package fxlogger

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	// Config is the configuration for the middleware.
	Config struct {
		// Logger is a custom instance of the logger to use.
		Logger *Logger
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper
		// BeforeNext is a function that is executed before the next handler is called.
		BeforeNext middleware.BeforeFunc
		// Enricher is a function that can be used to enrich the logger with additional information.
		Enricher Enricher
		// HandleError indicates whether to propagate errors up the middleware chain, so the global error handler can decide appropriate status code.
		HandleError bool
	}

	// Enricher is a function that can be used to enrich the logger with additional information.
	Enricher func(c echo.Context, logger zerolog.Context) zerolog.Context

	// Context is a wrapper around echo.Context that provides a logger.
	Context struct {
		echo.Context
		logger *Logger
	}
)

// NewContext returns a new Context.
func NewContext(ctx echo.Context, logger *Logger) *Context {
	return &Context{ctx, logger}
}

func (c *Context) Logger() echo.Logger {
	return c.logger
}

// Middleware returns a middleware which logs HTTP requests.
func Middleware(config Config) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}

	if config.Logger == nil {
		config.Logger = NewLogger(os.Stdout, WithTimestamp())
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			var err error
			req := c.Request()
			res := c.Response()
			start := time.Now()

			id := req.Header.Get(echo.HeaderXRequestID)

			if id == "" {
				id = res.Header().Get(echo.HeaderXRequestID)
			}

			cloned := false
			logger := config.Logger

			if id != "" {
				logger = FromLogger(logger.log, WithField(strings.ToLower(echo.HeaderXRequestID), id))
				cloned = true
			}

			if config.Enricher != nil {
				// to avoid mutation of shared instance
				if !cloned {
					logger = FromLogger(logger.log)
					cloned = true
				}

				logger.log = config.Enricher(c, logger.log.With()).Logger()
			}

			ctx := req.Context()

			if ctx == nil {
				ctx = context.Background()
			}

			// Pass logger down to request context
			c.SetRequest(req.WithContext(logger.WithContext(ctx)))
			c = NewContext(c, logger)

			if config.BeforeNext != nil {
				config.BeforeNext(c)
			}

			if err = next(c); err != nil {
				if config.HandleError {
					c.Error(err)
				}
			}

			stop := time.Now()

			var mainEvt *zerolog.Event
			if err != nil {
				mainEvt = logger.log.Err(err)
			} else {
				mainEvt = logger.log.Info()
			}

			mainEvt.Int("status", res.Status)
			mainEvt.Str("method", req.Method)
			mainEvt.Str("uri", req.RequestURI)
			mainEvt.Str("host", req.Host)
			mainEvt.Str("user_agent", req.UserAgent())
			mainEvt.Str("remote_ip", c.RealIP())
			mainEvt.Str("referer", req.Referer())
			mainEvt.Str("latency", stop.Sub(start).String())

			cl := req.Header.Get(echo.HeaderContentLength)
			if cl == "" {
				cl = "0"
			}

			mainEvt.Str("bytes_in", cl)
			mainEvt.Str("bytes_out", strconv.FormatInt(res.Size, 10))

			mainEvt.Send()

			return err
		}
	}
}
