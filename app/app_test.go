package app

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

var (
	testMux *http.ServeMux
)

func init() {
	// super sorry if you want to run these tests in CI, I'm lazy I know
	StaticDir, _ = filepath.Abs("../static")
	TemplateDir, _ = filepath.Abs("../tmpl")

	LoadTemplates()
}

func setup() {
	testMux = http.NewServeMux()
	testMux.Handle("/", Handler())
	// defined in posts.go
	apiClient = nil
}

func teardown() {
	apiClient = nil
	testMux = nil
}

func getHTML(t *testing.T, url *url.URL) (*goquery.Document, *httptest.ResponseRecorder) {
	// create the request
	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		t.Fatalf("Failed creating request %s", err)
	}

	// create test recorder
	rw := httptest.NewRecorder()
	rw.Body = new(bytes.Buffer)

	// ServeHTTP on the request and recorder
	testMux.ServeHTTP(rw, req)

	// use goquery NewDocumentFromReader from the recorders body
	doc, err := goquery.NewDocumentFromReader(rw.Body)
	if err != nil {
		t.Fatal(err)
	}

	// return the document and recorder
	return doc, rw
}
