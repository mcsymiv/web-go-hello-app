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

func main() {
	addr := flag.String("addr", ":8080", "Application port")
	env := flag.String("env", "uat", "Application environment")

	flag.Parse()

	db, err := run(env)
	if err != nil {
		log.Println("could not start the app")
		log.Println("possible solution: 'docker start postgres'")
		log.Fatal(err)
	}

	if db != nil {
		defer db.SQL.Close()
	}

	log.Println("Started app. Listen on port", *addr)
	srv := &http.Server{
		Addr:    *addr,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run(e *string) (*driver.DB, error) {

	// session type declaration
	gob.Register(models.Search{})
	gob.Register(models.User{})

	app.InProduction = false

	app.InfoLog = log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stdout, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Sesssion manager settings
	session = scs.New()

	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tmplCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Println("Can not create template from cache")
		return nil, err
	}

	app.UseCache = false
	app.TemplateCache = tmplCache

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	var repo *hand.Repository
	var db *driver.DB

	if *e == "prod" {
		// connect to postgres
		log.Println("connecting to db")
		db, err := driver.ConnectDB(fmt.Sprintf("host=postgres port=%s user=postgres dbname=db password=password", postgresPort))
		if err != nil {
			log.Fatal("cannot connect to db")
		}

		repo = hand.NewRepo(&app, db)

	} else if *e == "dev" {
		log.Println("local db")
		repo = hand.NewTestRepo(&app)

	} else if *e == "uat" {
		log.Println("connecting to local db")
		db, err = driver.ConnectDB(fmt.Sprintf("host=localhost port=%s user=postgres dbname=db password=password", postgresPort))
		if err != nil {
			log.Fatal("cannot connect to postgres container")
		}

		repo = hand.NewRepo(&app, db)

	} else {
		log.Fatal("invalid arg. Valid 'dev', 'uat' or no arguments")
	}

	hand.NewHandlers(repo)
	return db, nil
}
