package server

import (
	"github.com/ekkinox/fx-template/modules/fxgorm"
	"github.com/ekkinox/fx-template/modules/fxpubsub"

	"go.uber.org/fx"
)

func RegisterModules() fx.Option {
	return fx.Options(
		// orm
		fxgorm.FxGormModule,
		// pubsub
		fxpubsub.FxPubSubModule,
	)
}
