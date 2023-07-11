package worker

import (
	"context"

	"github.com/ekkinox/fx-template/internal/worker/pubsub"
	"github.com/ekkinox/fx-template/modules/fxbootstrapper"
	"go.uber.org/fx"
)

var WorkerBoostrapper = fxbootstrapper.NewBootstrapper().WithOptions(
	RegisterModules(),
	RegisterServices(),
)

func BootstrapWorker(ctx context.Context) *fx.App {
	return WorkerBoostrapper.BoostrapApp(
		fx.Invoke(func(w *pubsub.SubscribeWorker) {
			w.Run(ctx)
		}),
	)
}
