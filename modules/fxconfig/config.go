package fxconfig

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
}

const AppEnvProd = "prod"
const AppEnvDev = "dev"
const AppEnvTest = "test"

const DefaultAppEnv = AppEnvProd
const DefaultAppName = "my-app"
const DefaultAppVersion = "unknown"

func NewConfig() (*Config, error) {

	v := viper.New()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	v.SetConfigName(fmt.Sprintf("config.%s", FetchAppEnv()))
	if err := v.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	for _, key := range v.AllKeys() {
		val := v.GetString(key)
		if strings.Contains(val, "${") {
			v.Set(key, os.ExpandEnv(val))
		}
	}

	return &Config{v}, nil
}

func (c *Config) AppName() string {
	return c.GetString("app.name")
}

func (c *Config) AppEnv() string {
	return c.GetString("app.env")
}

func (c *Config) AppVersion() string {
	return c.GetString("app.version")
}

func (c *Config) AppDebug() bool {
	return c.GetBool("app.debug")
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", DefaultAppName)
	v.SetDefault("app.env", DefaultAppEnv)
	v.SetDefault("app.version", DefaultAppVersion)
	v.SetDefault("app.debug", false)
}
