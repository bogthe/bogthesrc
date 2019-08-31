package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/bogthe/bogthesrc/router"
	"github.com/gorilla/mux"
)

type handler func(w http.ResponseWriter, r *http.Request) error

type internalError struct {
	Message string `json:"message"`
}

func Handler() *mux.Router {
	r := router.API()
	r.Get(router.Post).Handler(handler(servePost))
	r.Get(router.PostCreate).Handler(handler(serveCreatePost))
	r.Get(router.Posts).Handler(handler(servePosts))

	return r
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			handleError(w, "Fatal", r)
		}
	}()

	err := h(w, r)
	if err != nil {
		handleError(w, err.Error(), nil)
	}
}

func handleError(w http.ResponseWriter, message string, v interface{}) {
	w.Header().Set("content-type", "application/json; charset=utf-8")

	code := http.StatusInternalServerError
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&internalError{Message: http.StatusText(code)})
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "[API:Error] %s\n", message)

	if v != nil {
		buf.Write(debug.Stack())
	}

	log.Println(buf.String())
}
