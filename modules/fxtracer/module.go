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
	fx.Invoke(func(*trace.TracerProvider) {}),
)

type FxTracerParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Factory   TracerProviderFactory
	Config    *fxconfig.Config
	Logger    *fxlogger.Logger
}

func NewFxTracerProvider(p FxTracerParam) (*trace.TracerProvider, error) {

	// exporter
	var exporter Exporter
	if p.Config.AppEnv() == fxconfig.Test {
		exporter = Memory
	} else {
		if p.Config.GetBool("modules.tracer.enabled") {
			exporter = FetchExporter(p.Config.GetString("modules.tracer.exporter"))
		} else {
			exporter = Noop
		}
	}

	tracerProvider, err := p.Factory.Create(
		WithName(p.Config.AppName()),
		WithExporter(exporter),
		WithCollector(p.Config.GetString("modules.tracer.collector")),
	)
	if err != nil {
		p.Logger.Error().Err(err).Msg("error creating tracer provider")

		return nil, err
	}

	p.LifeCycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if err = tracerProvider.ForceFlush(ctx); err != nil {
				p.Logger.Error().Err(err).Msg("error flushing tracer provider")

				return err
			}

			if exporter != Memory {
				if err = tracerProvider.Shutdown(ctx); err != nil {
					p.Logger.Error().Err(err).Msg("error while shutting down tracer provider")

					return err
				}
			}

			return nil
		},
	})

	return tracerProvider, nil
}
