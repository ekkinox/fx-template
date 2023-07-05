package worker

import (
	"context"

	"github.com/ekkinox/fx-template/modules/fxpubsub"

	"go.uber.org/fx"
)

func RegisterModules(ctx context.Context) fx.Option {
	return fx.Options(
		fxpubsub.FxPubSubModule,
	)
}
