package driver

import (
	"database/sql"
	"fmt"
	"github.com/porky256/course-project/internal/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type DB struct {
	DB *bun.DB
}

func ConnectSQL(config config.DBConfig) (*DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.User, config.Password, config.Host, config.Port, config.Name, config.SSLMode)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	err := sqldb.Ping()
	if err != nil {
		return nil, err
	}

	sqldb, err = setConfigs(config, sqldb)
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, pgdialect.New())
	return &DB{DB: db}, nil
}

func setConfigs(dbConfig config.DBConfig, sqlDB *sql.DB) (*sql.DB, error) {
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenDbConn)
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleDbConn)
	sqlDB.SetConnMaxLifetime(dbConfig.MaxDbLifetime)
	sqlDB.SetConnMaxIdleTime(dbConfig.MaxDbIdletime)
	return sqlDB, nil
}
