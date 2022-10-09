package hand

import (
	"net/http"
	"os"

	"github.com/mcsymiv/web-hello-world/pkg/render"
)

func Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func Home(w http.ResponseWriter, r *http.Request) {
	//render.RenderTemplate(w, "home.page.tmpl")
	render.RenderTemplate(w, "home.page.tmpl")
}

func About(w http.ResponseWriter, r *http.Request) {
	//render.RenderTemplate(w, "about.page.tmpl")
	render.RenderTemplate(w, "about.page.tmpl")
}

func Exit(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}
