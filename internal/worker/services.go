package worker

import (
	"context"

	"github.com/ekkinox/fx-template/internal/worker/pubsub"
	"go.uber.org/fx"
)

func RegisterServices(ctx context.Context) fx.Option {
	return fx.Provide(
		pubsub.NewSubscribeWorker,
	)
}
