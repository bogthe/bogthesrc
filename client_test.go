package bogthesrc

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient(nil)
	baseUrl, _ := url.Parse(server.URL)
	client.BaseUrl = baseUrl
}

func teardown() {
	server.Close()
}

func urlFor(t *testing.T, route string, values map[string]string) string {
	u, err := client.url(route, values, nil)
	if err != nil {
		t.Errorf("Failed url for route %s : %s", route, err)
	}

	return "/" + u.Path
}

func checkMethod(t *testing.T, r *http.Request, expected string) {
	if r.Method != expected {
		t.Errorf("Wrong http Method: w: %s g:%s", expected, r.Method)
	}
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		panic("Error encoding:" + err.Error())
	}
}
