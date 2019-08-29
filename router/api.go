package router

import "github.com/gorilla/mux"

func API() *mux.Router {
	r := mux.NewRouter()
	r.Path("/posts/{ID:.+}").Methods("GET").Name(Post)

	return r
}
