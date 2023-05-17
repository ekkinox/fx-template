package app

import (
	"github.com/ekkinox/fx-template/app/probes"
	"github.com/ekkinox/fx-template/app/services"
	"github.com/ekkinox/fx-template/modules/fxhealthchecker"
	"go.uber.org/fx"
)

func RegisterServices() fx.Option {
	return fx.Provide(
		services.NewTestService,
		fxhealthchecker.RegisterProbe(probes.NewTest1Probe),
		fxhealthchecker.RegisterProbe(probes.NewTest2Probe),
	)
}
