package dbrepo

import (
	"github.com/porky256/course-project/internal/config"
	"github.com/porky256/course-project/internal/repository"
	"github.com/uptrace/bun"
)

type postgressDB struct {
	App *config.AppConfig
	DB  *bun.DB
}

func NewPostgressDB(conn *bun.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgressDB{
		App: a,
		DB:  conn,
	}
}
