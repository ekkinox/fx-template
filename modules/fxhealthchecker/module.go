package fxhealthchecker

import (
	"go.uber.org/fx"
)

var FxHealthCheckerModule = fx.Module("health-checker",
	fx.Provide(
		NewFxHealthChecker,
	),
)

type FxHealthCheckerParam struct {
	fx.In
	Probes []HealthCheckerProbe `group:"health-checker-probes"`
}

func NewFxHealthChecker(p FxHealthCheckerParam) *HealthChecker {

	checker := NewHealthChecker()

	for _, probe := range p.Probes {
		checker.AddProbe(probe)
	}

	return checker
}
