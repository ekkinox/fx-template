package fxtracer

import (
	"context"

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
		exporter = FetchExporter(p.Config.GetString("tracing.exporter"))
	}

	tracerProvider, err := p.Factory.Create(
		WithName(p.Config.AppName()),
		WithExporter(exporter),
		WithCollector(p.Config.GetString("tracing.collector")),
	)
	if err != nil {
		p.Logger.Error().Err(err).Msg("error creating tracer provider")

		return nil, err
	}

	p.LifeCycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if err = tracerProvider.Shutdown(ctx); err != nil {
				p.Logger.Error().Err(err).Msg("error shutting down tracer provider")

				return err
			}

			return nil
		},
	})

	return tracerProvider, nil
}
