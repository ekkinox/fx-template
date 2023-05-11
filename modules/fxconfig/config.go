package fxconfig

import "github.com/spf13/viper"

type Config struct {
	AppName       string
	AppPort       int
	AppDebug      bool
	AppShouldFail bool
}

func NewConfig() *Config {

	viper.AutomaticEnv()
	setConfigDefaults()
	_ = viper.ReadInConfig()

	return &Config{
		AppName:       viper.GetString("APP_NAME"),
		AppPort:       viper.GetInt("APP_PORT"),
		AppDebug:      viper.GetBool("APP_DEBUG"),
		AppShouldFail: viper.GetBool("APP_SHOULD_FAIL"),
	}
}

func setConfigDefaults() {
	viper.SetDefault("APP_NAME", "my-app")
	viper.SetDefault("APP_PORT", 8080)
	viper.SetDefault("APP_DEBUG", true)
	viper.SetDefault("APP_SHOULD_FAIL", false)
}
