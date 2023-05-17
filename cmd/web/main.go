package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/porky256/course-project/internal/config"
	"github.com/porky256/course-project/internal/handlers"
	"github.com/porky256/course-project/internal/models"
	"github.com/porky256/course-project/internal/render"
	"log"
	"net/http"
	"time"
)

const port = ":8080"

// TODO move it to context
var app config.AppConfig

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
		return
	}

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

// run registers and initializes application
func run() error {
	gob.Register(models.Reservation{})

	cache, err := render.CreateTemplateCacheMap()

	//change it when production
	app.IsProduction = false

	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsProduction

	app.Session = session
	if err != nil {
		log.Fatal("can't create template cache: ", err)
		return err
	}
	app.TemplateCache = cache
	app.UseCache = false
	return nil
}
