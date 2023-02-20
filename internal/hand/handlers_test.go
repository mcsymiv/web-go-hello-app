package hand

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type testHandler struct {
	routes http.Handler
	server *httptest.Server
	client *http.Client
	paths  map[string]string
}

var h *testHandler

// handlerSetup
// Creates one testing server, client
// for all handler routes checks
func handlerSetup() *testHandler {
	r := getRoutes()
	s := httptest.NewTLSServer(r)
	c := s.Client()

	h = &testHandler{
		routes: r,
		server: s,
		client: c,
		paths: map[string]string{
			"index":    "/",
			"home":     "/home",
			"about":    "/about",
			"contact":  "/contact",
			"login":    "/user/login",
			"logout":   "/user/logout",
			"register": "/user/register",
		},
	}

	return h
}

// handlerTearDown
// references testHandler server for closing
func (h *testHandler) handlerTearDown() {
	h.server.Close()
}

// TestRoutePaths performs status check
// on GET handler routes
func TestRoutePaths(t *testing.T) {
	for _, p := range h.paths {
		res, err := h.client.Get(h.server.URL + p)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("%s endpoint failed with invalid status code. Expected: %d. But was: %d", p, http.StatusOK, res.StatusCode)
		}
	}
}

/*
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

func TestRepository_NoUserId_InSession_PostQuery(t *testing.T) {
	sq := "search_query=query"
	desc := "desc=description"
	rb := fmt.Sprintf("%s&%s", sq, desc)

	req, _ := http.NewRequest("POST", "/query", strings.NewReader(rb))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostQuery)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostQuery returns invalid code status on missing 'userId' in session, expected %d, but got %d", http.StatusTemporaryRedirect, rr.Code)
	}
}

func TestRepository_InvalidFormValue_CustomError_PostQuery(t *testing.T) {
	var u int = 1

	// missing search_query
	desc := "desc=description"

	req, _ := http.NewRequest("POST", "/query", strings.NewReader(desc))
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

func TestRepository_InserSearchToDB_Error_PostQuery(t *testing.T) {
	var u int = 1

	sq := "search_query=invalid"
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

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostQuery return invalid code status from DB on invalid 'search_query', expected %d, but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

}

func getCtx(r *http.Request) context.Context {
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println("unable to load session context from request", err)
	}

	return ctx
}
*/
