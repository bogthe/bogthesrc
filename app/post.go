package app

import (
	"net/http"

	"github.com/bogthe/bogthesrc"
)

var service = &bogthesrc.PostService{bogthesrc.NewClient(nil)}

func servePost(w http.ResponseWriter, r *http.Request) error {
	post, err := service.Get("1")
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
	posts, err := service.List(&bogthesrc.PostListOptions{})
	if err != nil {
		return err
	}

	return renderTemplate(w, r, PostListTemplate, http.StatusOK, struct {
		Posts []*bogthesrc.Post
	}{
		posts,
	})
}
