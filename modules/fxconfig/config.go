package fxconfig

import (
	"github.com/spf13/viper"
)

const (
	AppEnvProd = "prod"
	AppEnvDev  = "dev"
	AppEnvTest = "test"

	DefaultAppEnv     = AppEnvProd
	DefaultAppName    = "app"
	DefaultAppVersion = "unknown"
)

type Config struct {
	*viper.Viper
}

func (c *Config) AppName() string {
	return c.GetString("app.name")
}

func (c *Config) AppEnv() AppEnv {
	return FetchAppEnv(c.GetString("app.env"))
}

func (c *Config) AppVersion() string {
	return c.GetString("app.version")
}

func (c *Config) AppDebug() bool {
	return c.GetBool("app.debug")
}
