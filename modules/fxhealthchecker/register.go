package fxhealthchecker

import "go.uber.org/fx"

func AsHealthCheckerProbe(p any) any {
	return fx.Annotate(
		p,
		fx.As(new(HealthCheckerProbe)),
		fx.ResultTags(`group:"health-checker-probes"`),
	)
}
