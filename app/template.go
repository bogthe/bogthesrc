package app

import (
	"fmt"
	htmpl "html/template"
	"log"
	"net/http"
	"path/filepath"
)

const (
	HomeTemplate     = "home.html"
	PostTemplate     = "post/show.html"
	PostListTemplate = "post/list.html"
	ErrorTemplate    = "error/error.html"
)

var (
	TemplateDir string
)

var templates = make(map[string]*htmpl.Template)

func LoadTemplates() {
	err := parseTemplates([][]string{
		{HomeTemplate, "common.html", "layout.html"},
		{PostTemplate, "common.html", "layout.html"},
		{PostListTemplate, "common.html", "layout.html"},
		{ErrorTemplate, "common.html", "layout.html"},
	})

	if err != nil {
		log.Fatal(err)
	}
}

func renderTemplate(w http.ResponseWriter, r *http.Request, name string, status int, data interface{}) error {
	w.WriteHeader(status)
	// .Get is case insensitive
	if ct := w.Header().Get("content-type"); ct == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}

	tmpl := templates[name]
	if tmpl == nil {
		return fmt.Errorf("Couldn't find template for: %s", name)
	}

	return tmpl.Execute(w, data)
}

func parseTemplates(sets [][]string) error {
	for _, set := range sets {
		key := set[0]
		files := basePath(TemplateDir, set)
		tmpl, err := htmpl.New(key).ParseFiles(files...)
		if err != nil {
			return fmt.Errorf("Failed to parse template: %v, %s", set, err)
		}

		tmpl = tmpl.Lookup("ROOT")
		if tmpl == nil {
			return fmt.Errorf("No ROOT template found in set: %v", set)
		}

		templates[key] = tmpl
	}

	return nil
}

func basePath(base string, files []string) []string {
	paths := make([]string, len(files))
	for i := range files {
		paths[i] = filepath.Join(base, files[i])
	}

	return paths
}
