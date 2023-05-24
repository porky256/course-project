package config

import "time"

type DBConfig struct {
	MaxOpenDbConn int
	MaxIdleDbConn int
	MaxDbLifetime time.Duration
	MaxDbIdletime time.Duration
	Dsn           string
}
