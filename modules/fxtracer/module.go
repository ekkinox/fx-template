package fxtracer

import (
	"context"
	"fmt"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

var FxTracerModule = fx.Module(
	"tracer",
	fx.Provide(
		NewDefaultTracerProviderFactory,
		NewFxTracerProvider,
	),
)

type FxTracerParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Factory   TracerProviderFactory
	Config    *fxconfig.Config
	Logger    *fxlogger.Logger
}

func NewFxTracerProvider(p FxTracerParam) (*trace.TracerProvider, error) {

	exporter := Noop

	if p.Config.GetBool("tracing.enabled") {
		exporter = GetExporter(p.Config.GetString("tracing.exporter"))
	}

	fmt.Printf("******************************** config enqabled %s\n", p.Config.GetBool("tracing.enabled"))
	fmt.Printf("******************************** exporter config %s\n", p.Config.GetString("tracing.exporter"))
	fmt.Printf("******************************** exporter %s\n", exporter.String())

	tp, err := p.Factory.Create(
		WithName(p.Config.AppName()),
		WithExporter(exporter),
		WithCollector(p.Config.GetString("tracing.collector")),
	)
	if err != nil {
		return nil, err
	}

	p.LifeCycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if err = tp.Shutdown(ctx); err != nil {
				p.Logger.Error().Err(err).Msgf("error shutting down tracer provider: %v")
				return err
			}
			return nil
		},
	})

	return tp, nil
}
