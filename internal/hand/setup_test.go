package hand

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
	"github.com/mcsymiv/web-hello-world/internal/config"
	"github.com/mcsymiv/web-hello-world/internal/models"
	"github.com/mcsymiv/web-hello-world/internal/render"
)

var pathToTemplates string = "./../../templates"
var app config.AppConfig
var session *scs.SessionManager

func getRoutes() http.Handler {

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

	tmplCache, err := CreateTestTemplateCache()
	if err != nil {
		log.Println("Can not create template from cache")
	}

	app.TemplateCache = tmplCache
	app.UseCache = true

	repo := NewRepo(&app)
	NewHandlers(repo)

	render.NewRenderer(&app)

	var mux *chi.Mux = chi.NewRouter()

	// Middlewares
	// Gracefully absorb panics and prints the stack trace
	mux.Use(middleware.Recoverer)
	// CSRF

	// NoSurf not used for testing post requests '/home'
	// mux.Use(NoSurf)
	// Session
	mux.Use(SessionLoad)

	// get
	mux.Get("/", Repo.Index)
	mux.Get("/home", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/exit", Repo.Exit)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/result", Repo.SearchResult)

	// post
	mux.Post("/home", Repo.PostSearch)

	var fileServer http.Handler = http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

func NoSurf(next http.Handler) http.Handler {
	var csrfHandler *nosurf.CSRFHandler = nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// LoadAndSave loads and saves session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	tmplCache := map[string]*template.Template{}
	// get all files with match *.page.tmpl from ./templates folder
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return tmplCache, err
	}

	// range through all files with match *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		tmplSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return tmplCache, err
		}

		// look into layouts in ./template folder
		layoutFiles, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return tmplCache, err
		}

		if len(layoutFiles) > 0 {
			tmplSet, err = tmplSet.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return tmplCache, err
			}
		}

		tmplCache[name] = tmplSet
	}

	return tmplCache, nil
}
