package fxtracer

import (
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"go.uber.org/fx"
)

var FxTracerModule = fx.Module("tracer",
	fx.Provide(
		NewFxTracer,
	),
)

type FxTracerParam struct {
	fx.In
	Config *fxconfig.Config
}

type FxTracerResult struct {
	fx.Out
	TracerProvider *TracerProvider
}

func NewFxTracer(p FxTracerParam) FxTracerResult {
	return FxTracerResult{
		TracerProvider: NewTracerProvider(),
	}
}
