package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bogthe/bogthesrc"
	"github.com/bogthe/bogthesrc/datastore"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var (
	decoder = schema.NewDecoder()
	store   = &datastore.PostStore{datastore.NewDatastore(nil)}
)

func writeJSON(w http.ResponseWriter, v interface{}) error {
	data, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return err
	}

	w.Header().Set("content-type", "application/json; charset=utf-8")
	_, err = w.Write(data)

	return err
}

func servePost(w http.ResponseWriter, r *http.Request) error {
	// get ID from route
	id, ok := mux.Vars(r)["ID"]
	if !ok {
		return fmt.Errorf("No ID was specified")
	}

	post, err := store.Get(id)
	if err != nil {
		return err
	}

	// write to JSON
	return writeJSON(w, post)
}

func servePosts(w http.ResponseWriter, r *http.Request) error {
	var opt bogthesrc.PostListOptions
	if err := decoder.Decode(&opt, r.URL.Query()); err != nil {
		return err
	}

	posts, err := store.List(&opt)
	if err != nil {
		return err
	}

	return writeJSON(w, posts)
}
