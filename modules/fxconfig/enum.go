package fxconfig

import (
	"os"
	"strings"
)

type AppEnv int

const (
	Prod AppEnv = iota
	Dev
	Test
)

func (e AppEnv) String() string {
	switch e {
	case Prod:
		return AppEnvProd
	case Dev:
		return AppEnvDev
	case Test:
		return AppEnvTest
	default:
		return AppEnvProd
	}
}

func FetchAppEnvFromEnv() AppEnv {

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = DefaultAppEnv
	}

	return FetchAppEnv(env)
}

func FetchAppEnv(env string) AppEnv {

	switch strings.ToLower(env) {
	case AppEnvProd:
		return Prod
	case AppEnvDev:
		return Dev
	case AppEnvTest:
		return Test
	default:
		return Prod
	}
}
