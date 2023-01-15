package hand

import (
	"net/http"
	"time"

	"github.com/mcsymiv/web-hello-world/internal/config"
	"github.com/mcsymiv/web-hello-world/internal/driver"
	"github.com/mcsymiv/web-hello-world/internal/forms"
	"github.com/mcsymiv/web-hello-world/internal/models"
	"github.com/mcsymiv/web-hello-world/internal/render"
	"github.com/mcsymiv/web-hello-world/internal/repository"
	"github.com/mcsymiv/web-hello-world/internal/repository/dbrepo"
)

// Repo is the repository used by handlers
var Repo *Repository
var userId int

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

// NewTestRepo creates new repository with App config
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestDBRepo(a),
	}
}

// NewHandlers sets repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Index renders home page and puts user id into session
func (repo *Repository) Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// Home renders home MyQ page
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
		repo.App.ErrorLog.Println("can't parse form")
		repo.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	userId, ok := repo.App.Session.Get(r.Context(), "userId").(int)
	if !ok {
		repo.App.ErrorLog.Println("can not get 'userId' from session")
		repo.App.Session.Put(r.Context(), "error", "can not get 'userId' from Session")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	search := models.Search{
		Query:       r.Form.Get("search_query"),
		Link:        r.Form.Get("link"),
		Description: r.Form.Get("desc"),
		UserId:      userId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	form := forms.New(r.PostForm)

	form.Required("search_query", "desc")

	if !form.Valid() {
		data := make(map[string]interface{})

		data["invalid_search"] = search.Query
		data["invalid_desc"] = search.Description

		repo.App.ErrorLog.Println("invalid form values, check if 'search_query' or 'desc' are not empty")
		repo.App.Session.Put(r.Context(), "error", "invalid form value(s)")

		render.Template(w, r, "home.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})

		return
	}

	err = repo.DB.InsertSearch(search)
	if err != nil {
		repo.App.ErrorLog.Println("can't insert search to db")
		repo.App.Session.Put(r.Context(), "error", "can't insert search to db")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	repo.App.Session.Put(r.Context(), "new_query", &search)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// QueryResult renders query result page
func (repo *Repository) QueryResult(w http.ResponseWriter, r *http.Request) {
	userId, ok := repo.App.Session.Get(r.Context(), "userId").(int)
	if !ok {
		repo.App.ErrorLog.Println("Can not get 'userId' from session")
		repo.App.Session.Put(r.Context(), "error", "Can not get 'userId' from Session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		return
	}

	query := r.URL.Query().Get("query")

	form := forms.New(r.URL.Query())
	form.Required("query")

	if !form.Valid() {
		repo.App.ErrorLog.Println("invalid form values, check if 'query' is empty")
		repo.App.Session.Put(r.Context(), "error", "invalid form value")

		render.Template(w, r, "home.page.tmpl", &models.TemplateData{
			Form: form,
		})

		return
	}

	search, err := repo.DB.GetUserSearchesByUserIdAndPartialTextQuery(userId, query)
	if err != nil {
		repo.App.ErrorLog.Println("unable to get searches from DB")
		repo.App.ErrorLog.Println("unable to get search for user from DB on partial text query", err)
		http.Redirect(w, r, "/result", http.StatusTemporaryRedirect)

		return
	}

	data := make(map[string]interface{})
	data["search"] = search

	stringMap := make(map[string]string)
	stringMap["searchQuery"] = query

	render.Template(w, r, "result.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// Login renders login page
func (repo *Repository) Login(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// Logout, terminates user session, redirects to home page
func (repo *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = repo.App.Session.Destroy(r.Context())
	_ = repo.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// PostLogin handles user login auth logic
func (repo *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	// good practice to renew session token on login
	_ = repo.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		repo.App.ErrorLog.Println("unable to parse form", err)
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")

	if !form.Valid() {
		repo.App.ErrorLog.Println("required fields are empty, check 'email' or 'password'")
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})

		return
	}

	email := r.Form.Get("email")
	pass := r.Form.Get("password")

	id, _, err := repo.DB.AuthenticateUser(email, pass)
	if err != nil {
		repo.App.ErrorLog.Println("unable to authenticate user", err)
		repo.App.Session.Put(r.Context(), "error", "invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)

		return
	}

	repo.App.InfoLog.Println("successfully authenticated")
	repo.App.Session.Put(r.Context(), "userId", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

// About renders about app page
func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Contact renders contact page
func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// Dashboard renders all user queries to manage
func (repo *Repository) Dashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "myq-dashboard.page.tmpl", &models.TemplateData{})
}

// MyqSearches renders searches
func (repo *Repository) MyqSearches(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "myq-searches.page.tmpl", &models.TemplateData{})
}

// MyqUsers shows total users count
func (repo *Repository) MyqUsers(w http.ResponseWriter, r *http.Request) {
	user_count, err := repo.DB.GetUsersCount()
	if err != nil {
		repo.App.ErrorLog.Println("unable to get user count from DB", err)
		http.Redirect(w, r, "/myq/dashboard", http.StatusTemporaryRedirect)

		return
	}

	data := make(map[string]interface{})
	data["users"] = user_count

	render.Template(w, r, "myq-users.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
