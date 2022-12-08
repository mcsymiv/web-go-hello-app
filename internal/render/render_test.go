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

func TestTemplate(t *testing.T) {
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

	err = Template(&rw, r, "query.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
	if err != nil {
		t.Errorf("Error writing template to browser: %v", err)
	}

	err = Template(&rw, r, "non-existing.page.tmpl", &models.TemplateData{})
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

func TestNewRenderer(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Errorf("Can not create template cache, %v", err)
	}
}
