package config

import "time"

type DBConfig struct {
	MaxOpenDbConn int
	MaxIdleDbConn int
	MaxDbLifetime time.Duration
	MaxDbIdletime time.Duration
	Host          string
	Port          string
	Name          string
	User          string
	Password      string
	SSLMode       string
}
