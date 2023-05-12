package fxconfig

import "github.com/spf13/viper"

type Config struct {
	AppConfig AppConfig
	*viper.Viper
}

type AppConfig struct {
	Name  string
	Port  int
	Debug bool
}

func NewConfig() *Config {

	v := viper.New()

	v.AutomaticEnv()
	setAppConfigDefaults(v)
	_ = viper.ReadInConfig()

	return &Config{
		AppConfig{
			Name:  v.GetString("APP_NAME"),
			Port:  v.GetInt("APP_PORT"),
			Debug: v.GetBool("APP_DEBUG"),
		},
		v,
	}
}

func setAppConfigDefaults(v *viper.Viper) {
	v.SetDefault("APP_NAME", "my-app")
	v.SetDefault("APP_PORT", 8080)
	v.SetDefault("APP_DEBUG", false)
}
