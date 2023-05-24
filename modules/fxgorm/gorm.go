package fxgorm

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGorm(driver string, dsn string) (*gorm.DB, error) {

	var dial gorm.Dialector
	switch GetDriver(driver) {
	case Sqlite:
		dial = sqlite.Open(dsn)
	case Mysql:
		dial = mysql.Open(dsn)
	case Postgres:
		dial = postgres.Open(dsn)
	}

	return gorm.Open(dial, &gorm.Config{})
}
