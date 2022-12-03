package hand

import (
	"net/http"
	"os"

	"github.com/mcsymiv/web-hello-world/internal/config"
	"github.com/mcsymiv/web-hello-world/internal/driver"
	"github.com/mcsymiv/web-hello-world/internal/forms"
	"github.com/mcsymiv/web-hello-world/internal/helpers"
	"github.com/mcsymiv/web-hello-world/internal/models"
	"github.com/mcsymiv/web-hello-world/internal/render"
	"github.com/mcsymiv/web-hello-world/internal/repository"
	"github.com/mcsymiv/web-hello-world/internal/repository/dbrepo"
)

// Repo is the repository used by handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates new repository with App config
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(a, db.SQL),
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
	var emptySearch models.Search
	data := make(map[string]interface{})
	data["search"] = emptySearch

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (repo *Repository) PostHome(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	search := models.Search{
		Query:     r.Form.Get("query"),
		StartDate: r.Form.Get("start"),
		EndDate:   r.Form.Get("end"),
	}

	form := forms.New(r.PostForm)

	form.Required("query", "start", "end")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["search"] = search

		render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	repo.App.Session.Put(r.Context(), "search", search)
	http.Redirect(w, r, "/result", http.StatusSeeOther)
}

func (repo *Repository) SearchResult(w http.ResponseWriter, r *http.Request) {
	search, ok := repo.App.Session.Get(r.Context(), "search").(models.Search)
	if !ok {
		repo.App.ErrorLog.Println("Can not get 'Search' data from Session")
		repo.App.Session.Put(r.Context(), "error", "Can not get 'Search result' data from Session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// remove 'search' from session
	repo.App.Session.Remove(r.Context(), "search")

	data := make(map[string]interface{})
	data["search"] = search

	render.RenderTemplate(w, r, "result.page.tmpl", &models.TemplateData{
		Data: data,
	})
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
