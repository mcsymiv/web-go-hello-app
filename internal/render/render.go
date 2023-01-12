package render

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/mcsymiv/web-hello-world/internal/config"
	"github.com/mcsymiv/web-hello-world/internal/models"
)

var pathToTemplates string = "./templates"
var app *config.AppConfig

// sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

// AddDefaultTemplateData adds data that can be used across all templates
func AddDefaultTemplateData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "userId") {
		td.IsAuth = 1
	}
	return td
}

func Template(w http.ResponseWriter, r *http.Request, tmpl string, tmplData *models.TemplateData) error {
	var tmplCache map[string]*template.Template

	if app.UseCache {
		// get template cache from app config
		tmplCache = app.TemplateCache
	} else {
		tmplCache, _ = CreateTemplateCache()
	}

	// get template
	t, ok := tmplCache[tmpl]
	if !ok {
		return errors.New("Can not get template from cache")
	}

	buf := new(bytes.Buffer)

	// Add default data to template data
	tmplData = AddDefaultTemplateData(tmplData, r)

	// Render template
	err := t.Execute(buf, tmplData)
	if err != nil {
		log.Fatalf("Can not execute template cache render: %v", err)
	}

	// render template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Printf("Can not write to response writer buffer: %v", err)
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	tmplCache := map[string]*template.Template{}

	// get all files with match *.page.tmpl from ./templates folder
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return tmplCache, err
	}

	// range through all files with match *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		tmplSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return tmplCache, err
		}

		// look into layouts in ./template folder
		layoutFiles, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return tmplCache, err
		}

		if len(layoutFiles) > 0 {
			tmplSet, err = tmplSet.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return tmplCache, err
			}
		}

		tmplCache[name] = tmplSet
	}

	return tmplCache, nil
}
