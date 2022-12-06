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

var timeFormat string = "2006-01-01"

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

func (repo *Repository) PostSearch(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	userId := 1

	repo.App.Session.Put(r.Context(), "user_id", userId)
	repo.App.Session.Put(r.Context(), "query", r.Form.Get("query"))

	search := models.Search{
		Query:  r.Form.Get("query"),
		UserId: userId,
	}

	form := forms.New(r.PostForm)

	form.Required("query")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["search"] = search

		render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	err = repo.DB.InsertSearch(search)
	if err != nil {
		helpers.ServerError(w, err)
	}

	http.Redirect(w, r, "/result", http.StatusSeeOther)
}

func (repo *Repository) SearchResult(w http.ResponseWriter, r *http.Request) {
	userId, ok := repo.App.Session.Get(r.Context(), "user_id").(int)
	if !ok {
		repo.App.ErrorLog.Println("Can not get 'user_id' from session")
		repo.App.Session.Put(r.Context(), "error", "Can not get 'user_id' from Session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	query, ok := repo.App.Session.Get(r.Context(), "query").(string)
	if !ok {
		repo.App.ErrorLog.Println("Can not get 'query' from session")
		repo.App.Session.Put(r.Context(), "query_error", "Can not get 'query' from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// remove 'user_id' from session
	repo.App.Session.Remove(r.Context(), "user_id")
	repo.App.Session.Remove(r.Context(), "query")

	search, err := repo.DB.GetUserSearch(query, userId)
	if err != nil {
		helpers.ServerError(w, err)
	}

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
