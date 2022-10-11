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

	mux.Get("/", hand.Repo.Index)
	mux.Get("/home", hand.Repo.Home)
	mux.Get("/about", hand.Repo.About)
	mux.Get("/exit", hand.Repo.Exit)

	return mux
}
