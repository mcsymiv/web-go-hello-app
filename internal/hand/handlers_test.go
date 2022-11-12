package hand

import (
	"net/http"
	"net/http/httptest"
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

		}
	}

}
