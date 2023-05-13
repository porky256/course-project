package main

import (
	"fmt"
	"github.com/porky256/course-project/pkg/config"
	"github.com/porky256/course-project/pkg/handlers"
	"github.com/porky256/course-project/pkg/render"
	"log"
	"net/http"
)

const port = ":8080"

func main() {
	var app config.AppConfig
	cache, err := render.CreateTemplateCacheMap()
	if err != nil {
		log.Fatal("can't create template cache: ", err)
	}
	app.TemplateCache = cache
	app.UseCache = false

	newRender := render.NewRender(&app)
	newHandler := handlers.NewHandlers(&app, newRender)

	server := http.Server{
		Addr:    port,
		Handler: routes(&app, newHandler),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("starting application on port", port)
}
