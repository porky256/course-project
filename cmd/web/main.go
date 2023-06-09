package main

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
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

const host = "localhost:8080"
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

	app.InfoLog.Println("Connecting to DB...")
	db, err := driver.ConnectSQL(dbconfig)
	if err != nil {
		app.ErrorLog.Fatal(err)
		return
	}
	app.InfoLog.Println("Connection established")
	defer db.DB.Close()

	defer close(app.MailChan)
	listenForEmail()

	newRender := render.NewRender(&app)
	newHandler := handlers.NewHandlers(&app, newRender, db)

	server := http.Server{
		Addr:    host,
		Handler: routes(&app, newHandler),
	}

	app.InfoLog.Println("starting application on port", port)
	err = server.ListenAndServe()
	if err != nil {
		app.ErrorLog.Fatal(err)
	}
	time.Date(1, 2, 3, 4, 5, 6, 7, time.UTC)
}

// run registers and initializes application
func run() error {
	gob.Register(models.Reservation{})
	gob.Register(models.Room{})
	gob.Register(models.User{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})
	gob.Register(map[string]int{})

	godotenv.Load()
	dbUser, dbPassword, dbName, dbHost, dbPort, dbSSLMode, inProduction, useCache :=
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_SSLMODE"),
		os.Getenv("IN_PRODUCTION"),
		os.Getenv("USE_CACHE")

	dbconfig = config.DBConfig{
		User:          dbUser,
		Password:      dbPassword,
		Name:          dbName,
		Host:          dbHost,
		Port:          dbPort,
		SSLMode:       dbSSLMode,
		MaxOpenDbConn: 10,
		MaxIdleDbConn: 5,
		MaxDbLifetime: 24 * time.Hour,
		MaxDbIdletime: 24 * time.Hour,
	}

	app.RootPath = "./"
	cache, err := render.CreateTemplateCacheMap(&app)
	if err != nil {
		app.ErrorLog.Fatal("can't create template cache: ", err)
		return err
	}
	//change it when production
	app.IsProduction = inProduction == "true"
	app.UseCache = useCache == "true"

	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsProduction

	app.TemplateCache = cache
	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.DateLayout = "2006-01-02"
	app.MailChan = make(chan models.MailData)

	app.Session = session

	helpers.NewHelpers(&app)
	return nil
}
