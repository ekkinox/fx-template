package fxgorm

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

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

	orm, err := f.createDatabase(appliedOpts.Driver.String(), appliedOpts.Dsn, appliedOpts.Config)
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

func (f *DefaultGormFactory) createDatabase(driver string, dsn string, config gorm.Config) (*gorm.DB, error) {

	var dial gorm.Dialector
	switch FetchDriver(driver) {
	case Sqlite3:
		dial = sqlite.Open(dsn)
	case Mysql:
		dial = mysql.Open(dsn)
	case Postgres:
		dial = postgres.Open(dsn)
	case SqlServer:
		dial = sqlserver.Open(dsn)
	}

	return gorm.Open(dial, &config)
}
