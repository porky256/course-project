package handlers

import (
	"github.com/porky256/course-project/internal/config"
	"github.com/porky256/course-project/internal/driver"
	"github.com/porky256/course-project/internal/render"
	"github.com/porky256/course-project/internal/repository"
	"github.com/porky256/course-project/internal/repository/dbrepo"
	mock_dbrepo "github.com/porky256/course-project/internal/repository/mock"
)

type Handlers struct {
	app    *config.AppConfig
	render *render.Render
	DB     repository.DatabaseRepo
}

func NewHandlers(app *config.AppConfig, render *render.Render, db *driver.DB) *Handlers {
	return &Handlers{
		app:    app,
		render: render,
		DB:     dbrepo.NewPostgresDB(db.DB, app),
	}
}

func NewTestHandlers(app *config.AppConfig, render *render.Render, db *mock_dbrepo.MockDatabaseRepo) *Handlers {
	return &Handlers{
		app:    app,
		render: render,
		DB:     db,
	}
}
