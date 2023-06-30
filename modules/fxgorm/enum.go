package fxgorm

import (
	"strings"

	"gorm.io/gorm/logger"
)

type Driver int

const (
	Unknown Driver = iota
	Sqlite
	Mysql
	Postgres
	SqlServer
)

func (d Driver) String() string {
	switch d {
	case Sqlite:
		return "sqlite"
	case Mysql:
		return "mysql"
	case Postgres:
		return "postgres"
	case SqlServer:
		return "sqlserver"
	default:
		return "unknown"
	}
}

func FetchDriver(driver string) Driver {
	switch strings.ToLower(driver) {
	case "sqlite":
		return Sqlite
	case "mysql":
		return Mysql
	case "postgres":
		return Postgres
	case "sqlserver":
		return SqlServer
	default:
		return Unknown
	}
}

func FetchLogLevel(level string) logger.LogLevel {
	switch strings.ToLower(level) {
	case "silent":
		return logger.Silent
	case "info":
		return logger.Info
	case "warn":
		return logger.Warn
	case "error":
		return logger.Error
	default:
		return logger.Silent
	}
}
