package fxgorm

import "gorm.io/gorm"

type GormFactory interface {
	Create(options ...GormOption) (*gorm.DB, error)
}

type DefaultGormFactory struct{}

func NewDefaultGormFactory() GormFactory {
	return &DefaultGormFactory{}
}

func (f *DefaultGormFactory) Create(options ...GormOption) (*gorm.DB, error) {

	appliedOpts := defaultGormOptions
	for _, applyOpt := range options {
		applyOpt(&appliedOpts)
	}

	orm, err := NewGorm(appliedOpts.Driver.String(), appliedOpts.Dsn, appliedOpts.Config)
	if err != nil {
		return nil, err
	}

	for _, plugin := range appliedOpts.Plugins {
		err = orm.Use(plugin)
		if err != nil {
			return nil, err
		}
	}

	return orm, nil
}
