package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mcsymiv/web-hello-world/pkg/config"
	"github.com/mcsymiv/web-hello-world/pkg/hand"
)

func routes(app *config.AppConfig) http.Handler {
	var mux *chi.Mux = chi.NewRouter()

	// Middlewares
	// Gracefully absorb panics and prints the stack trace
	mux.Use(middleware.Recoverer)
	// CSRF
	mux.Use(NoSurf)
	// Session
	mux.Use(SessionLoad)

	// get
	mux.Get("/", hand.Repo.Index)
	mux.Get("/home", hand.Repo.Home)
	mux.Get("/about", hand.Repo.About)
	mux.Get("/exit", hand.Repo.Exit)
	mux.Get("/contact", hand.Repo.Contact)

	// post
	mux.Post("/home", hand.Repo.PostHome)

	// ajax
	mux.Get("/home-json", hand.Repo.HomeJson)

	var fileServer http.Handler = http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
