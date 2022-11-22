package render

import (
	"net/http"
	"testing"

	"github.com/mcsymiv/web-hello-world/internal/forms"
	"github.com/mcsymiv/web-hello-world/internal/models"
)

func TestAddDefaultTemplateData(t *testing.T) {
	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")

	result := AddDefaultTemplateData(&td, r)
	if result.Flash != "123" {
		t.Errorf("failed 'AddDefaultTemplateData'. Flash value expected: '123', but was %s", result.Flash)
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var rw testRenderResponseWriter
	var data map[string]interface{}

	err = RenderTemplate(&rw, r, "home.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
	if err != nil {
		t.Errorf("Error writing template to browser: %v", err)
	}

	err = RenderTemplate(&rw, r, "non-existing.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Errorf("Got non-existing template: %v", err)
	}
}

func getSession() (*http.Request, error) {

	r, err := http.NewRequest("GET", "/some-req", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Errorf("Can not create template cache, %v", err)
	}
}
