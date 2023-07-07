package fxconfig

import (
	"os"

	"go.uber.org/fx"
)

var FxConfigModule = fx.Module("config",
	fx.Provide(
		NewDefaultConfigFactory,
		NewFxConfig,
	),
)

type FxConfigParam struct {
	fx.In
	Factory ConfigFactory
}

func NewFxConfig(p FxConfigParam) (*Config, error) {
	return p.Factory.Create(
		WithFileName("config"),
		WithFilePaths(
			".",
			"./configs",
			os.Getenv("APP_CONFIG_PATH"),
		),
	)
}
