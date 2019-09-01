package app

import (
	"net/http"
	"strconv"

	"github.com/bogthe/bogthesrc"
	"github.com/gorilla/mux"
)

var apiClient = bogthesrc.NewApiClient(nil)

func servePost(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["ID"])
	if err != nil {
		return err
	}

	post, err := apiClient.Posts.Get(id)
	if err != nil {
		return err
	}

	return renderTemplate(w, r, PostTemplate, http.StatusOK, struct {
		Post *bogthesrc.Post
	}{
		post,
	})
}

func servePosts(w http.ResponseWriter, r *http.Request) error {
	posts, err := apiClient.Posts.List(&bogthesrc.PostListOptions{})
	if err != nil {
		return err
	}

	return renderTemplate(w, r, PostListTemplate, http.StatusOK, struct {
		Posts []*bogthesrc.Post
	}{
		posts,
	})
}
