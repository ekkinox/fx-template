package fxgorm

import (
	"strings"
)

type Driver int

const (
	Unknown Driver = iota
	Sqlite
	Mysql
	Postgres
)

func (d Driver) String() string {
	switch d {
	case Sqlite:
		return "sqlite"
	case Mysql:
		return "mysql"
	case Postgres:
		return "postgres"
	default:
		return "unknown"
	}
}

func GetDriver(driver string) Driver {
	switch strings.ToLower(driver) {
	case "sqlite":
		return Sqlite
	case "mysql":
		return Mysql
	case "postgres":
		return Postgres
	default:
		return Unknown
	}
}
