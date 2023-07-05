package worker

import (
	"context"

	"github.com/ekkinox/fx-template/internal/worker/pubsub"
	"go.uber.org/fx"
)

func RegisterOverrides(ctx context.Context) fx.Option {
	return fx.Options(
		// pubsub subscriber invocation
		fx.Invoke(func(w *pubsub.SubscribeWorker) {
			w.Run(ctx)
		}),
	)
}
