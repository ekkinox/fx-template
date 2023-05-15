package fxhttpserver

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const tracerKey = "otel-go-contrib-tracer-labstack-echo"

func GetCtxLogger(ctx echo.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx.Request().Context())
}

func GetCtxTracer(ctx echo.Context) trace.Tracer {
	t := ctx.Get(tracerKey)
	if t != nil {
		return ctx.Get(tracerKey).(trace.Tracer)
	}

	return otel.Tracer("default")
}
