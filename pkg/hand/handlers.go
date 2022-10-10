package hand

import (
	"net/http"
	"os"

	"github.com/mcsymiv/web-hello-world/pkg/config"
	"github.com/mcsymiv/web-hello-world/pkg/render"
)

// Repo is the repository used by handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates new repository with App config
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (repo *Repository) Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	//render.RenderTemplate(w, "home.page.tmpl")
	render.RenderTemplate(w, "home.page.tmpl")
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	//render.RenderTemplate(w, "about.page.tmpl")
	render.RenderTemplate(w, "about.page.tmpl")
}

func (repo *Repository) Exit(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}
