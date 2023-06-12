package server

import (
	"github.com/ekkinox/fx-template/internal/server/repository"
	"github.com/ekkinox/fx-template/modules/fxgorm"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"go.uber.org/fx"
)

func RegisterServices() fx.Option {
	return fx.Provide(
		// health check probes
		fxhealthchecker.AsProbe(fxgorm.NewGormProbe),
		// repositories
		repository.NewPostRepository,
	)
}
