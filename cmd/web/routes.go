package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/porky256/course-project/pkg/config"
	"github.com/porky256/course-project/pkg/handlers"
	"net/http"
)

func routes(app *config.AppConfig, handler *handlers.Handlers) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", http.HandlerFunc(handler.Home))
	mux.Get("/about", http.HandlerFunc(handler.About))

	return mux
}
