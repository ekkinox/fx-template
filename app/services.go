package app

import (
	"github.com/ekkinox/fx-template/app/repository"
	"github.com/ekkinox/fx-template/modules/fxgorm"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"go.uber.org/fx"
)

func RegisterServices() fx.Option {
	return fx.Provide(
		// health check probes
		fxhealthchecker.RegisterProbe(fxgorm.NewGormProbe),
		// repositories
		repository.NewPostRepository,
	)
}
