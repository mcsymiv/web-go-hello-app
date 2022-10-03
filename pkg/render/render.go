package render

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/mcsymiv/web-hello-world/pkg/er"
)

var templateCacheMap = make(map[string]*template.Template)

func RenderTemplate(w http.ResponseWriter, tmpl string) {
	p, _ := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl")
	err := p.Execute(w, nil)
	er.CheckErr("error parsing template", err)
}

func RenderCachedTemplate(w http.ResponseWriter, t string) {
	var tmpl *template.Template
	var err error

	_, inMap := templateCacheMap[t]
	if !inMap {
		fmt.Println("Creating template and adding to cache")
		err = createCachedTemplate(t)
		if err != nil {
			fmt.Printf("Unable create template in cache: %v", err)
		}
	} else {
		fmt.Println("Using cached template")
	}

	tmpl = templateCacheMap[t]

	err = tmpl.Execute(w, nil)
	if err != nil {
		fmt.Printf("Unable to execute template parsing: %v", err)
	}
}

func createCachedTemplate(t string) error {
	templates := []string{
		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.tmpl",
	}

	tmpl, err := template.ParseFiles(templates...)
	if err != nil {
		fmt.Printf("Unable to parse templates: %v", err)
		return err
	}

	templateCacheMap[t] = tmpl

	return nil
}
