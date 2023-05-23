package fxconfig

import "go.uber.org/fx"

var FxConfigModule = fx.Module("config",
	fx.Provide(
		NewFxConfig,
	),
)

func NewFxConfig() (*Config, error) {
	return NewConfig()
}
