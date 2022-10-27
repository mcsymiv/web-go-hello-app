package hand

import (
	"log"
	"net/http"
	"os"

	"github.com/mcsymiv/web-hello-world/internal/config"
	"github.com/mcsymiv/web-hello-world/internal/forms"
	"github.com/mcsymiv/web-hello-world/internal/models"
	"github.com/mcsymiv/web-hello-world/internal/render"
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
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (repo *Repository) PostHome(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Error on parse form from /home-post", err)
		return
	}

	search := models.Search{
		Query:     r.Form.Get("query"),
		StartDate: r.Form.Get("start_date"),
		EndDate:   r.Form.Get("end_date"),
	}

	form := forms.New(r.PostForm)

	form.HasRequired("query", r)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["search"] = search

		render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) Exit(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}
