package fxlogger

import (
	"context"
	"github.com/labstack/echo/v4"

	"github.com/rs/zerolog"
)

// WithContext returns a new context with the provided logger.
func (l Logger) WithContext(ctx context.Context) context.Context {
	return l.Unwrap().WithContext(ctx)
}

// Ctx returns a logger from the provided context.
// If no logger is found in the context, a new one is created.
func Ctx(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}

// EchoCtx returns a logger from the provided echo context.
// If no logger is found in the context, a new one is created.
func EchoCtx(ctx echo.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx.Request().Context())
}
