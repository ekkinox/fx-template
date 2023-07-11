package fxgrpcserver

import (
	"context"

	"github.com/ekkinox/fx-template/modules/fxlogger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func CtxLogger(ctx context.Context) *fxlogger.Logger {
	return fxlogger.CtxLogger(ctx)
}

func CtxTracer(ctx context.Context) trace.Tracer {
	return otel.Tracer("fxgrpcserver-tracer")
}
