package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/mcsymiv/web-hello-world/internal/config"
	"github.com/mcsymiv/web-hello-world/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	// session type declaration
	gob.Register(models.Search{})

	testApp.InProduction = false
	testApp.InfoLog = log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	testApp.ErrorLog = log.New(os.Stdout, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Sesssion manager settings
	session = scs.New()

	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session
	app = &testApp

	os.Exit(m.Run())
}

type testRenderResponseWriter struct{}

func (rw *testRenderResponseWriter) Header() http.Header {
	var h http.Header
	return h
}

func (rw *testRenderResponseWriter) WriteHeader(i int) {}

func (rw *testRenderResponseWriter) Write(b []byte) (int, error) {
	l := len(b)
	return l, nil
}
