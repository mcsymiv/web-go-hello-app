package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/mcsymiv/web-hello-world/pkg/config"
)

var app *config.AppConfig

// sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func RenderTemplate(w http.ResponseWriter, tmpl string) {
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
		log.Fatal("Could not get template from tmplCache")
	}

	buf := new(bytes.Buffer)

	_ = t.Execute(buf, nil)

	// render template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	tmplCache := map[string]*template.Template{}
	// get all files with match *.page.tmpl from ./templates folder
	pages, err := filepath.Glob("./templates/*.page.tmpl")
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
		layoutFiles, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return tmplCache, err
		}

		if len(layoutFiles) > 0 {
			tmplSet, err = tmplSet.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return tmplCache, err
			}
		}

		tmplCache[name] = tmplSet
	}

	return tmplCache, nil
}
