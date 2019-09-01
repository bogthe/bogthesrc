package router

import "github.com/gorilla/mux"

const (
	PostCreateForm = "post:create-form"
)

func App() *mux.Router {
	r := mux.NewRouter()
	r.Path("/").Methods("GET").Name(Home)
	r.Path("/posts").Methods("GET").Name(Posts)
	r.Path("/submit").Methods("GET").Name(PostCreateForm)
	r.Path("/posts").Methods("POST").Name(PostCreate)
	r.Path("/posts/{ID:.+}").Methods("GET").Name(Post)

	return r
}
