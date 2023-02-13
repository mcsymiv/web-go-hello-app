package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mcsymiv/web-hello-world/internal/config"
	"github.com/mcsymiv/web-hello-world/internal/hand"
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
	// Method convertion
	mux.Use(Method)

	// get
	mux.Get("/", hand.Repo.Index)
	mux.Get("/home", hand.Repo.Home)
	//	mux.Get("/query", hand.Repo.Query)
	mux.Get("/about", hand.Repo.About)
	mux.Get("/contact", hand.Repo.Contact)
	mux.Get("/result", hand.Repo.QueryResult)
	mux.Get("/user/login", hand.Repo.Login)
	mux.Get("/user/logout", hand.Repo.Logout)
	mux.Get("/user/register", hand.Repo.Register)

	// post
	mux.Post("/query", hand.Repo.PostQuery)
	mux.Post("/user/login", hand.Repo.PostLogin)
	mux.Post("/user/register", hand.Repo.PostRegister)

	var fileServer http.Handler = http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/myq", func(mux chi.Router) {
		if app.InProduction {
			mux.Use(Authenticated)
		}

		mux.Get("/dashboard", hand.Repo.Dashboard)
		mux.Get("/searches", hand.Repo.MyqSearches)
		mux.Get("/searches/{id}/edit", hand.Repo.MyqSearchesView)
		mux.Get("/searches/{id}/delete", hand.Repo.MyqSearchDelete)

		mux.Put("/searches/{id}/edit", hand.Repo.MyqSearchEdit)
	})

	return mux
}
