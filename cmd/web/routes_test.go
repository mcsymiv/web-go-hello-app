package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi"
	"github.com/mcsymiv/web-hello-world/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// passed
	default:
		t.Error(fmt.Sprintf("Failed 'routes()'. Returned type not *chi.Mux, was: %t", v))
	}
}
