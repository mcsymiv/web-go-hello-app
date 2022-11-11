package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

type noSurfHandler struct{}

func (sh *noSurfHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
