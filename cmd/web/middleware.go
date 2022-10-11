package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// NoSurf adds csrf protection to all POST requests
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
