package fxgrpcserver

import (
	"context"

	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/ekkinox/fx-template/modules/fxtracer"
)

func CtxLogger(ctx context.Context) *fxlogger.Logger {
	return fxlogger.CtxLogger(ctx)
}

func CtxTracer(ctx context.Context) *fxtracer.Tracer {
	return fxtracer.CtxTracer(ctx)
}
