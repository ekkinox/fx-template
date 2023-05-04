package fxconfig

import "go.uber.org/fx"

var ConfigModule = fx.Module("config",
	fx.Provide(
		NewConfig,
	),
)

type ConfigResult struct {
	fx.Out
	Config *Config
}

func NewConfig() ConfigResult {
	return ConfigResult{
		Config: newConfig(),
	}
}
