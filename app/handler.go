package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	ReloadTemplates bool
	StaticDir       string
)

func Handler() *mux.Router {
	router := mux.NewRouter()
	router.Path("/").Methods("GET").Handler(handler(serveHome))

	return router
}

type handler func(w http.ResponseWriter, r *http.Request) error

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(w, r)
	if err != nil {
		log.Fatal("This should not have happened")
	}
}
