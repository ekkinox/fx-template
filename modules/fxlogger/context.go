package fxlogger

import (
	"context"
	"github.com/labstack/echo/v4"

	"github.com/rs/zerolog"
)

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
