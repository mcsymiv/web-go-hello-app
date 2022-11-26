package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/mcsymiv/web-hello-world/internal/config"
	"github.com/mcsymiv/web-hello-world/internal/hand"
	"github.com/mcsymiv/web-hello-world/internal/helpers"
	"github.com/mcsymiv/web-hello-world/internal/models"
	"github.com/mcsymiv/web-hello-world/internal/render"
)

const port = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	err := run()
	if err != nil {
		log.Println("Could not start the app")
		log.Fatal(err)
	}

	fmt.Println("Started app. Listen on port :8080")
	srv := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	// session type declaration
	gob.Register(models.Search{})

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
		return nil
	}

	app.UseCache = false
	app.TemplateCache = tmplCache

	repo := hand.NewRepo(&app)
	hand.NewHandlers(repo)

	render.NewTemplates(&app)
	helpers.NewHelpers(&app)

	return nil
}
