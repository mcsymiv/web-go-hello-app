package render

import (
	"net/http"
	"testing"

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
