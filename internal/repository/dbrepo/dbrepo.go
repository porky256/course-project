package dbrepo

import (
	"github.com/porky256/course-project/internal/config"
	"github.com/porky256/course-project/internal/repository"
	"github.com/uptrace/bun"
)

type postgresDB struct {
	App *config.AppConfig
	DB  *bun.DB
}

// NewPostgresDB creates a new postgres DB entity
func NewPostgresDB(conn *bun.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDB{
		App: a,
		DB:  conn,
	}
}
