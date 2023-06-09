package worker

import (
	"github.com/ekkinox/fx-template/modules/fxpubsub"

	"go.uber.org/fx"
)

func RegisterModules() fx.Option {
	return fx.Options(
		// pubsub
		fxpubsub.FxPubSubModule,
	)
}
