package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/mcsymiv/web-hello-world/internal/config"
)

var app *config.AppConfig

func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// IsAuthenticated checks if user is logged in
func IsAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "userId")
}
