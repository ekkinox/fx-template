package fxtracer

import (
	"context"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

var FxTracerModule = fx.Module(
	"tracer",
	fx.Provide(
		NewFxTracer,
	),
	fx.Decorate(func(e *echo.Echo) *echo.Echo {
		e.Use(otelecho.Middleware("my-app"))
		e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Logger().Info("tracer middleware active")
				return next(c)
			}
		})

		return e
	}),
	fx.Invoke(func(*trace.TracerProvider) {}),
)

type FxTracerParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Config    *fxconfig.Config
	Logger    *fxlogger.Logger
}

func NewFxTracer(p FxTracerParam) (*trace.TracerProvider, error) {

	tp, err := NewTracerProvider(p.Config, p.Logger)
	if err != nil {
		return nil, err
	}

	p.LifeCycle.Append(fx.Hook{
		// stop
		OnStop: func(ctx context.Context) error {
			if err := tp.Shutdown(ctx); err != nil {
				p.Logger.Errorf("error shutting down tracer provider: %v", err)
				return err
			}
			return nil
		},
	})

	return tp, nil
}
