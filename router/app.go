package router

import "github.com/gorilla/mux"

func App() *mux.Router {
	r := mux.NewRouter()
	r.Path("/").Methods("GET").Name(Home)
	r.Path("/posts").Methods("GET").Name(Posts)
	r.Path("/posts/{ID:.+}").Methods("GET").Name(Post)

	return r
}
