package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/porky256/course-project/internal/config"
	"github.com/porky256/course-project/internal/driver"
	"github.com/porky256/course-project/internal/handlers"
	"github.com/porky256/course-project/internal/helpers"
	"github.com/porky256/course-project/internal/models"
	"github.com/porky256/course-project/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

const port = ":8080"

// TODO move it to context
var app config.AppConfig
var dbconfig config.DBConfig

func main() {
	err := run()
	if err != nil {
		app.ErrorLog.Fatal(err)
		return
	}
	fmt.Println("Connecting to DB...")
	db, err := driver.ConnectSQL(dbconfig)
	if err != nil {
		app.ErrorLog.Fatal(err)
		return
	}
	fmt.Println("Connection established")
	defer db.DB.Close()
	newRender := render.NewRender(&app)
	newHandler := handlers.NewHandlers(&app, newRender, db)

	server := http.Server{
		Addr:    port,
		Handler: routes(&app, newHandler),
	}

	fmt.Println("starting application on port", port)
	err = server.ListenAndServe()
	if err != nil {
		app.ErrorLog.Fatal(err)
	}
}

// run registers and initializes application
func run() error {
	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.User{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})
	dbconfig = config.DBConfig{
		Dsn:           "postgres://postgres:2341@0.0.0.0:5432/db?sslmode=disable",
		MaxOpenDbConn: 10,
		MaxIdleDbConn: 5,
		MaxDbLifetime: 24 * time.Hour,
		MaxDbIdletime: 24 * time.Hour,
	}

	app.RootPath = "./"
	cache, err := render.CreateTemplateCacheMap(&app)

	//change it when production
	app.IsProduction = false

	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsProduction

	app.TemplateCache = cache
	app.UseCache = false
	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.DateLayout = "2006-01-02"

	app.Session = session
	if err != nil {
		app.ErrorLog.Fatal("can't create template cache: ", err)
		return err
	}

	helpers.NewHelpers(&app)
	return nil
}
