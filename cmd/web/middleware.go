package main

import (
	"net/http"
	"strings"

	"github.com/justinas/nosurf"
	"github.com/mcsymiv/web-hello-world/internal/helpers"
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

// Authenticated checks if user is logged in
func Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Log in first!")
			app.InfoLog.Println("unable to auth user in middleware")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)

			return
		}

		next.ServeHTTP(w, r)
	})
}

// Method middleware is provide by @logrusorgru from #355 issue go-chi/chi
// Package PPD implements HTTP middleware to replace reqeust method by
// POST-form "_method" value. It supports PUT, PATCH and DELETE methods.
// The Method is the middleware.
func Method(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			switch strings.ToLower(r.PostFormValue("_method")) {
			case "put":
				r.Method = http.MethodPut
			case "patch":
				r.Method = http.MethodPatch
			case "delete":
				r.Method = http.MethodDelete
			default:
			}
		}
		next.ServeHTTP(w, r)
	})
}
