package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/mcsymiv/web-hello-world/internal/config"
	"github.com/mcsymiv/web-hello-world/internal/driver"
	"github.com/mcsymiv/web-hello-world/internal/hand"
	"github.com/mcsymiv/web-hello-world/internal/helpers"
	"github.com/mcsymiv/web-hello-world/internal/models"
	"github.com/mcsymiv/web-hello-world/internal/render"
)

const postgresPort = "5432"

var app config.AppConfig
var session *scs.SessionManager
var srv *http.Server
var db *driver.DB
var repo *hand.Repository

func main() {
	err := run()
	if err != nil {
		app.InfoLog.Println("could not start the app")
		app.InfoLog.Println("possible solution: 'docker start postgres'")
		app.ErrorLog.Fatal(err)
	}
}

func run() error {

	// session type declaration
	gob.Register(models.Search{})
	gob.Register(models.User{})

	// app flags
	addr := flag.String("addr", ":8080", "Application port")
	dbname := flag.String("dbname", "db", "DB name")
	dbport := flag.String("dbport", "5432", "DB port")
	dbpass := flag.String("dbpass", "", "DB password")
	dbhost := flag.String("dbhost", "localhost", "DB host")
	dbuser := flag.String("dbuser", "postgres", "DB user")
	app.InProduction = *flag.Bool("prod", false, "In production")
	app.UseCache = *flag.Bool("cache", false, "Use cache")

	flag.Parse()

	// app logging
	app.InfoLog = log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = log.New(os.Stdout, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Sesssion manager settings
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// app templates
	tmplCache, err := render.CreateTemplateCache()
	if err != nil {
		app.InfoLog.Println("Can not create template from cache")
		return err
	}

	app.TemplateCache = tmplCache

	// connect to postgres
	app.InfoLog.Println("connecting to db")
	db, err := driver.ConnectDB(fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", *dbhost, *dbport, *dbuser, *dbname, *dbpass))
	if err != nil {
		app.ErrorLog.Fatal("cannot connect to db")

		return err
	}

	defer db.SQL.Close()

	repo = hand.NewRepo(&app, db)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)
	hand.NewHandlers(repo)

	app.InfoLog.Println(fmt.Sprintf("serving on port: %s", *addr))
	srv = &http.Server{
		Addr:     *addr,
		ErrorLog: app.ErrorLog,
		Handler:  routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		app.InfoLog.Fatal(err)
		return err
	}

	return nil
}
