package fxgorm

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func NewGorm(driver string, dsn string, config gorm.Config) (*gorm.DB, error) {

	var dial gorm.Dialector
	switch FetchDriver(driver) {
	case Sqlite:
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
