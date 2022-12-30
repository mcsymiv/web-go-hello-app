package hand

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type postData struct {
	key   string
	value string
}

var testGetHandlers = []struct {
	name       string
	url        string
	method     string
	statusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	// creates test server for the lifetime of the test
	testServer := httptest.NewTLSServer(routes)

	// closes server after test is done
	defer testServer.Close()

	for _, test := range testGetHandlers {
		res, err := testServer.Client().Get(testServer.URL + test.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if res.StatusCode != test.statusCode {
			t.Errorf("%s failed with invalid status code. Expected: %d. But was: %d", test.name, test.statusCode, res.StatusCode)
		}
	}
}

func TestRepository_Success_QueryResult(t *testing.T) {

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

func TestRepository_Redirect_Without_UserId_And_Query_In_Session_QueryResult(t *testing.T) {
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

func TestRepository_Success_WithoutQueryInSession_QueryResult(t *testing.T) {
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

func TestRepository_Redirect_Error_From_DB(t *testing.T) {
	// user_id must be equal to 5 to return error from test_repo on GetUserSearchesByUserIdAndPartialTextQuery DB query
	var u int = 5

	req, _ := http.NewRequest("GET", "/result", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	session.Put(ctx, "user_id", u)

	h := http.HandlerFunc(Repo.QueryResult)

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("failed get /result, expected status %d, but was: %d", http.StatusTemporaryRedirect, rr.Code)
	}
}

func TestRepository_PostQuery(t *testing.T) {
	var u int = 1

	sq := "search_query=query"
	desc := "desc=description"
	rb := fmt.Sprintf("%s&%s", sq, desc)

	req, _ := http.NewRequest("POST", "/query", strings.NewReader(rb))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	session.Put(ctx, "user_id", u)

	handler := http.HandlerFunc(Repo.PostQuery)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostQuery return invalid code status, expected %d, but got %d", http.StatusSeeOther, rr.Code)
	}
}

func TestRepository_ParseFormError_PostQuery(t *testing.T) {
	req, _ := http.NewRequest("POST", "/query", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostQuery)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostQuery without body returns invalid code status, expected %d, but got %d", http.StatusTemporaryRedirect, rr.Code)
	}
}

func getCtx(r *http.Request) context.Context {
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println("unable to load session context from request", err)
	}

	return ctx
}
