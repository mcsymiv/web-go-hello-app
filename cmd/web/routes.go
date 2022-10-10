package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/mcsymiv/web-hello-world/pkg/config"
	"github.com/mcsymiv/web-hello-world/pkg/hand"
)

func routes(app *config.AppConfig) http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(hand.Repo.Index))
	mux.Get("/home", http.HandlerFunc(hand.Repo.Home))
	mux.Get("/about", http.HandlerFunc(hand.Repo.About))
	mux.Get("/exit", http.HandlerFunc(hand.Repo.Exit))

	return mux
}
