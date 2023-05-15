package fxhttpserver

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func GetCtxLogger(ctx echo.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx.Request().Context())
}

func GetCtxTracer(ctx echo.Context) trace.Tracer {
	t := ctx.Get("otel-go-contrib-tracer-labstack-echo")
	if t != nil {
		return ctx.Get("otel-go-contrib-tracer-labstack-echo").(trace.Tracer)
	}

	return otel.Tracer("default")
}
