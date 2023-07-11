package fxbootstrapper

import (
	"testing"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"github.com/ekkinox/fx-template/modules/fxtracer"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type Bootstrapper struct {
	defaultOptions []fx.Option
}

func NewBootstrapper() *Bootstrapper {
	return &Bootstrapper{
		defaultOptions: []fx.Option{
			fxconfig.FxConfigModule,
			fxlogger.FxLoggerModule,
			fxtracer.FxTracerModule,
			fxhealthchecker.FxHealthCheckerModule,
		},
	}
}

func (b *Bootstrapper) WithOptions(option ...fx.Option) *Bootstrapper {
	b.defaultOptions = append(b.defaultOptions, option...)

	return b
}

func (b *Bootstrapper) BoostrapApp(bootstrapOptions ...fx.Option) *fx.App {
	return fx.New(
		fx.WithLogger(fxlogger.FxEventLogger),
		fx.Options(b.defaultOptions...),
		fx.Options(bootstrapOptions...),
	)
}

func (b *Bootstrapper) BoostrapAndRunApp(bootstrapOptions ...fx.Option) {
	b.BoostrapApp(bootstrapOptions...).Run()
}

func (b *Bootstrapper) BoostrapTestApp(t testing.TB, bootstrapOptions ...fx.Option) *fxtest.App {

	t.Setenv("APP_ENV", "test")

	return fxtest.New(
		t,
		fx.NopLogger,
		fx.Options(b.defaultOptions...),
		fx.Options(bootstrapOptions...),
	)
}

func (b *Bootstrapper) BoostrapAndRunTestApp(t testing.TB, bootstrapOptions ...fx.Option) {
	b.BoostrapTestApp(t, bootstrapOptions...).RequireStart().RequireStop()
}
