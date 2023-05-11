package app

import (
	"github.com/ekkinox/fx-template/app/services"
	"go.uber.org/fx"
)

func RegisterServices() fx.Option {
	return fx.Provide(
		services.NewTestService,
	)
}
