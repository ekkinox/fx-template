package server

import (
	"context"

	"github.com/ekkinox/fx-template/internal/repository"
	"github.com/ekkinox/fx-template/modules/fxgorm"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"github.com/ekkinox/fx-template/modules/fxpubsub"
	"go.uber.org/fx"
)

func RegisterServices(ctx context.Context) fx.Option {
	return fx.Provide(
		// health check probes
		fxhealthchecker.AsHealthCheckerProbe(fxgorm.NewGormProbe),
		fxhealthchecker.AsHealthCheckerProbe(fxpubsub.NewPubSubProbe),
		// repositories
		repository.NewPostRepository,
	)
}
