package fxconfig

import "go.uber.org/fx"

var FxConfigModule = fx.Module("config",
	fx.Provide(
		NewFxConfig,
	),
)

type FxConfigResult struct {
	fx.Out
	Config *Config
}

func NewFxConfig() FxConfigResult {
	return FxConfigResult{
		Config: NewConfig(),
	}
}
