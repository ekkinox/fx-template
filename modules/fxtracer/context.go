package fxtracer

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const ctxKey = "otel-go-contrib-tracer-labstack-echo"

func CtxTracer(ctx context.Context) *Tracer {

	ct := ctx.Value(ctxKey)

	var t trace.Tracer
	if ct != nil {
		t = ctx.Value(ctxKey).(trace.Tracer)
	} else {
		t = otel.Tracer("default")
	}

	return NewTracer(t, ctx)
}
