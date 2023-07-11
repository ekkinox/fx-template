package worker

import (
	"github.com/ekkinox/fx-template/internal/worker/pubsub"
	"go.uber.org/fx"
)

func RegisterServices() fx.Option {
	return fx.Provide(
		pubsub.NewSubscribeWorker,
	)
}
