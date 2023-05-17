package fxhttpserver

import (
	"github.com/ekkinox/fx-template/modules/fxtracer"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const tracerKey = "otel-go-contrib-tracer-labstack-echo"

func GetCtxLogger(ctx echo.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx.Request().Context())
}

func GetCtxTracer(ctx echo.Context) *fxtracer.Tracer {

	ct := ctx.Get(tracerKey)

	var t trace.Tracer
	if ct != nil {
		t = ctx.Get(tracerKey).(trace.Tracer)
	} else {
		t = otel.Tracer("default")
	}

	return fxtracer.NewTracer(t, ctx.Request().Context())
}
