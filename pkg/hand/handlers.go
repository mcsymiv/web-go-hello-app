package hand

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mcsymiv/web-hello-world/pkg/config"
	"github.com/mcsymiv/web-hello-world/pkg/models"
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
	var remoteIp string = r.RemoteAddr
	repo.App.Session.Put(r.Context(), "remoteIp", remoteIp)

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) PostHome(w http.ResponseWriter, r *http.Request) {
	var startDate string = r.FormValue("start")
	var endDate string = r.FormValue("end")
	w.Write([]byte(fmt.Sprintf("Posted to home route with data: %s, %s", startDate, endDate)))
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)

	stringMap["remoteIp"] = repo.App.Session.GetString(r.Context(), "remoteIp")

	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) Exit(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}

func (repo *Repository) HomeJson(w http.ResponseWriter, r *http.Request) {
	homeJson := struct {
		Ok      bool   `json:"ok"`
		Message string `json:"msg"`
	}{
		Ok:      true,
		Message: "Message from json",
	}

	res, err := json.MarshalIndent(homeJson, "", "\t")
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
