package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func RenderTemplate(w http.ResponseWriter, tmpl string) {
	// create template cache
	tmplCache, err := createTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	// get template
	t, ok := tmplCache[tmpl]
	if !ok {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)

	err = t.Execute(buf, nil)
	if err != nil {
		log.Println(err)
	}

	// render template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func createTemplateCache() (map[string]*template.Template, error) {
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
