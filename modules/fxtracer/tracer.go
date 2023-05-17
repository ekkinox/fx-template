package fxtracer

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	tracer  trace.Tracer
	context context.Context
}

func NewTracer(tracer trace.Tracer, context context.Context) *Tracer {
	return &Tracer{
		tracer:  tracer,
		context: context,
	}
}

func (t *Tracer) Start(spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return t.tracer.Start(t.context, spanName, opts...)
}
