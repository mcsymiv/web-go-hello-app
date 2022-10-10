package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mcsymiv/web-hello-world/pkg/config"
	"github.com/mcsymiv/web-hello-world/pkg/hand"
	"github.com/mcsymiv/web-hello-world/pkg/render"
)

const port = ":8080"

func main() {
	var app config.AppConfig

	tmplCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Can not create template from cache")
	}

	app.UseCache = false
	app.TemplateCache = tmplCache

	repo := hand.NewRepo(&app)
	hand.NewHandlers(repo)

	render.NewTemplates(&app)

	http.HandleFunc("/", hand.Repo.Index)
	http.HandleFunc("/home", hand.Repo.Home)
	http.HandleFunc("/about", hand.Repo.About)
	http.HandleFunc("/exit", hand.Repo.Exit)

	fmt.Println("Started app.\nListen on port :8080")
	_ = http.ListenAndServe(port, nil)
}
