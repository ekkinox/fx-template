package fxgorm

import (
	"context"
	"github.com/ekkinox/fx-template/modules/fxconfig"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var FxGormModule = fx.Module(
	"gorm",
	fx.Provide(
		NewFxGorm,
	),
	fx.Invoke(func(*gorm.DB) {}),
)

type FxGormParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Config    *fxconfig.Config
}

func NewFxGorm(p FxGormParam) (*gorm.DB, error) {

	// orm
	orm, err := NewGorm(p.Config.GetString("database.driver"), p.Config.GetString("database.dsn"))
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
