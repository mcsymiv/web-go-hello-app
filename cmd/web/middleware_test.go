package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {

	var noSurf noSurfHandler

	h := NoSurf(&noSurf)

	switch v := h.(type) {
	case http.Handler:
		// passed
	default:
		t.Error(fmt.Sprintf("Failed 'NoSurf'. Returned type is not http.Handler, %t", v))
	}
}

func TestSessionLoad(t *testing.T) {

	var noSurf noSurfHandler

	h := SessionLoad(&noSurf)

	switch v := h.(type) {
	case http.Handler:
		// passed
	default:
		t.Error(fmt.Sprintf("Failed 'SessionLoad'. Returned type is not http.Handler, %t", v))
	}
}
