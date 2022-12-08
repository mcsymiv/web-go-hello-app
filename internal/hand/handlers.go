package hand

import (
	"net/http"
	"os"
	"time"

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

// Index renders home page and puts user id into session
func (repo *Repository) Index(w http.ResponseWriter, r *http.Request) {
	var userId int = 1
	repo.App.Session.Put(r.Context(), "user_id", userId)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// Home renders home MyQ form
func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	var emptySearch models.Search

	data := make(map[string]interface{})
	data["search"] = emptySearch

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostQuery inserts new query into database with validation fields
func (repo *Repository) PostQuery(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	repo.App.Session.Put(r.Context(), "query", r.Form.Get("query"))

	userId := repo.App.Session.Get(r.Context(), "user_id").(int)

	search := models.Search{
		Query:       r.Form.Get("query"),
		Link:        r.Form.Get("link"),
		Description: r.Form.Get("desc"),
		UserId:      userId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	form := forms.New(r.PostForm)

	form.Required("query", "desc")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["search"] = search

		render.Template(w, r, "home.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	repo.App.InfoLog.Println("successfully get query form values", search)

	err = repo.DB.InsertSearch(search)
	if err != nil {
		helpers.ServerError(w, err)
	}

	http.Redirect(w, r, "/result", http.StatusSeeOther)
}

// QueryResult renders query result page from session
func (repo *Repository) QueryResult(w http.ResponseWriter, r *http.Request) {
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
		// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		// return
		repo.App.InfoLog.Println("setting query from url")
		query = r.URL.Query().Get("query")
	}

	search, err := repo.DB.GetUserSearchesByUserIdAndPartialTextQuery(userId, query)
	if err != nil {
		helpers.ServerError(w, err)
	}

	repo.App.InfoLog.Println("resutl from db", search)

	data := make(map[string]interface{})
	data["search"] = search

	stringMap := make(map[string]string)
	stringMap["sessionQuery"] = query

	// remove 'query' from session
	repo.App.Session.Remove(r.Context(), "query")

	render.Template(w, r, "result.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// About renders about app page
func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Contact renders contact page
func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// Exit kills app, removes user_id from session
// Todo: remove after graceful shutdown implementation
func (repo *Repository) Exit(w http.ResponseWriter, r *http.Request) {
	repo.App.Session.Remove(r.Context(), "user_id")
	os.Exit(0)
}
