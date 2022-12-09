package hand

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var tests = []struct {
	name       string
	url        string
	method     string
	params     []postData
	statusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	// creates test server for the lifetime of the test
	testServer := httptest.NewTLSServer(routes)

	// closes server after test is done
	defer testServer.Close()

	for _, test := range tests {
		if test.method == "GET" {
			res, err := testServer.Client().Get(testServer.URL + test.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if res.StatusCode != test.statusCode {
				t.Errorf("%s failed with invalid status code. Expected: %d. But was: %d", test.name, test.statusCode, res.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, v := range test.params {
				values.Add(v.key, v.value)
			}

			res, err := testServer.Client().PostForm(testServer.URL+test.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if res.StatusCode != test.statusCode {
				t.Errorf("%s failed with invalid status code. Expected: %d. But was: %d", test.name, test.statusCode, res.StatusCode)
			}
		}
	}
}

func TestRepository_SuccessSearch(t *testing.T) {

	// put in session 'user_id' and 'query'
	var u int = 1
	var q string = "test query"

	req, _ := http.NewRequest("GET", "/result", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	session.Put(ctx, "user_id", u)
	session.Put(ctx, "query", q)

	h := http.HandlerFunc(Repo.QueryResult)

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("failed get /result, expected status %d, but was: %d", http.StatusOK, rr.Code)
	}
}

func TestRepository_RedirectSearch(t *testing.T) {
	req, _ := http.NewRequest("GET", "/result", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	h := http.HandlerFunc(Repo.QueryResult)

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("failed get /result without 'user_id' and 'query' in session, expected status %d, but was: %d", http.StatusTemporaryRedirect, rr.Code)
	}
}

func TestRepository_RedirectSearch_WithoutQuery(t *testing.T) {
	var u int = 1

	req, _ := http.NewRequest("GET", "/result", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	session.Put(ctx, "user_id", u)

	h := http.HandlerFunc(Repo.QueryResult)

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("failed get /result without 'user_id' and 'query' in session, expected status %d, but was: %d", http.StatusOK, rr.Code)
	}
}

func getCtx(r *http.Request) context.Context {
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println("unable to load session context from request", err)
	}

	return ctx
}
