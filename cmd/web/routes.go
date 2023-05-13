package main

import (
	"github.com/bmizerany/pat"
	"github.com/porky256/course-project/pkg/config"
	"github.com/porky256/course-project/pkg/handlers"
	"net/http"
)

func routes(app *config.AppConfig, handler *handlers.Handlers) http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handler.Home))
	mux.Get("/about", http.HandlerFunc(handler.About))

	return mux
}
