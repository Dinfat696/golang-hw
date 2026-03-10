package database

import (
	"time"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func ConnectDbWithCfg(dbDriverName string, dsn string) *sqlx.DB {
	DB = sqlx.MustConnect(dbDriverName, dsn)

	DB.SetMaxIdleConns(5)
	DB.SetMaxOpenConns(20)
	DB.SetConnMaxLifetime(1 * time.Minute)
	DB.SetConnMaxIdleTime(10 * time.Minute)
	return DB
}
