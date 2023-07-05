package server

import (
	"context"

	"github.com/ekkinox/fx-template/modules/fxgorm"
	"github.com/ekkinox/fx-template/modules/fxpubsub"

	"go.uber.org/fx"
)

func RegisterModules(ctx context.Context) fx.Option {
	return fx.Options(
		// orm
		fxgorm.FxGormModule,
		// pubsub
		fxpubsub.FxPubSubModule,
	)
}
