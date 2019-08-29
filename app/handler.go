package app

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
)

var (
	ReloadTemplates bool
	StaticDir       string
)

func Handler() *mux.Router {
	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(StaticDir))))
	router.Path("/").Methods("GET").Handler(handler(serveHome))

	return router
}

type handler func(w http.ResponseWriter, r *http.Request) error

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ReloadTemplates {
		loadTemplates()
	}

	runHandler(w, r, h)
}

func runHandler(w http.ResponseWriter, r *http.Request, fn handler) {
	defer func() {
		if rv := recover(); rv != nil {
			err := errors.New("runHandler error")
			logError(r, err, rv)
			handleError(w, r, http.StatusInternalServerError, err)
		}
	}()

	err := fn(w, r)
	if err != nil {
		logError(r, err, nil)
		handleError(w, r, http.StatusInternalServerError, err)
	}
}

func logError(r *http.Request, err error, rv interface{}) {
	if err != nil {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "Error serving the route %s: %s\n", mux.CurrentRoute(r).GetName(), err)
		if rv != nil {
			fmt.Fprint(&buf, rv)
			buf.Write(debug.Stack())
		}

		log.Print(buf.String())
	}
}

func handleError(w http.ResponseWriter, r *http.Request, status int, err error) {
	renderErr := renderTemplate(w, r, ErrorTemplate, status, &struct {
		StatusCode int
		Status     string
		Error      error
	}{
		StatusCode: status,
		Status:     http.StatusText(status),
		Error:      err,
	})

	if renderErr != nil {
		log.Fatalf("Failed to render error template: %s", renderErr)
	}
}
