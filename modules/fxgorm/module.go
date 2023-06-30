package fxgorm

import (
	"context"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var FxGormModule = fx.Module(
	"gorm",
	fx.Provide(
		NewDefaultGormFactory,
		NewFxGorm,
	),
)

type FxGormParam struct {
	fx.In
	LifeCycle      fx.Lifecycle
	Factory        GormFactory
	Config         *fxconfig.Config
	Logger         *fxlogger.Logger
	TracerProvider *trace.TracerProvider
}

func NewFxGorm(p FxGormParam) (*gorm.DB, error) {

	config := gorm.Config{
		Logger: NewGormLogger(
			p.Logger,
			p.Config.GetBool("gorm.logger.with_values"),
		).LogMode(
			FetchLogLevel(p.Config.GetString("gorm.logger.level")),
		),
	}

	var plugins []gorm.Plugin
	if p.Config.GetBool("gorm.tracer.enabled") {
		plugins = append(
			plugins,
			NewGormTracerPlugin(p.TracerProvider, p.Config.GetBool("gorm.tracer.with_values")),
		)
	}

	orm, err := p.Factory.Create(
		WithDsn(p.Config.GetString("gorm.dsn")),
		WithDriver(FetchDriver(p.Config.GetString("gorm.driver"))),
		WithConfig(config),
		WithPlugins(plugins...),
	)

	if err != nil {
		return nil, err
	}

	// lifecycle
	p.LifeCycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			db, err := orm.DB()
			if err != nil {
				return err
			}

			return db.Close()
		},
	})

	return orm, nil
}
