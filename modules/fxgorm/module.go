package fxgorm

import (
	"context"

	"github.com/ekkinox/fx-template/modules/fxconfig"
	"github.com/ekkinox/fx-template/modules/fxlogger"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var FxGormModule = fx.Module(
	"gorm",
	fx.Provide(
		NewFxGorm,
	),
)

type FxGormParam struct {
	fx.In
	LifeCycle      fx.Lifecycle
	Config         *fxconfig.Config
	Logger         *fxlogger.Logger
	TracerProvider *trace.TracerProvider
}

func NewFxGorm(p FxGormParam) (*gorm.DB, error) {

	// orm
	logLevel := logger.Error
	if p.Config.AppDebug() {
		logLevel = logger.Info
	}

	config := gorm.Config{
		Logger: NewGormLogger(p.Logger).LogMode(logLevel),
	}

	orm, err := NewGorm(p.Config.GetString("database.driver"), p.Config.GetString("database.dsn"), config)
	if err != nil {
		return nil, err
	}

	if p.Config.AppDebug() {
		err = orm.Use(NewGormTracerPlugin(p.TracerProvider))
		if err != nil {
			return nil, err
		}
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
