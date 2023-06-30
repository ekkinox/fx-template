package fxgorm

import (
	"gorm.io/gorm"
)

type options struct {
	Dsn     string
	Driver  Driver
	Config  gorm.Config
	Plugins []gorm.Plugin
}

var defaultGormOptions = options{
	Dsn:     "",
	Driver:  Unknown,
	Config:  gorm.Config{},
	Plugins: []gorm.Plugin{},
}

type GormOption func(o *options)

func WithDsn(d string) GormOption {
	return func(o *options) {
		o.Dsn = d
	}
}

func WithDriver(d Driver) GormOption {
	return func(o *options) {
		o.Driver = d
	}
}

func WithConfig(c gorm.Config) GormOption {
	return func(o *options) {
		o.Config = c
	}
}

func WithPlugins(p ...gorm.Plugin) GormOption {
	return func(o *options) {
		o.Plugins = p
	}
}
