package bogthesrc

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/bogthe/bogthesrc/router"
)

func TestGETPost(t *testing.T) {
	setup()
	defer teardown()

	wantPost := &Post{ID: "1"}
	wantPosts := []*Post{wantPost}
	calledPost := false
	calledPosts := false

	mux.HandleFunc(urlFor(t, router.Post, map[string]string{"ID": "1"}), func(w http.ResponseWriter, r *http.Request) {
		calledPost = true

		checkMethod(t, r, "GET")
		writeJSON(w, wantPost)
	})

	mux.HandleFunc(urlFor(t, router.Posts, nil), func(w http.ResponseWriter, r *http.Request) {
		calledPosts = true

		checkMethod(t, r, "GET")
		writeJSON(w, wantPosts)
	})

	posts := &PostService{Client: client}

	t.Run("Can get a post by id", func(t *testing.T) {
		actualPost, err := posts.Get("1")

		if err != nil {
			t.Errorf("Failed to get post :%s", err)
		}

		if !calledPost {
			t.Error("Handler func not called")
		}

		if !reflect.DeepEqual(wantPost, actualPost) {
			t.Errorf("Posts aren't the same want: %+v, got: %+v", wantPost, actualPost)
		}
	})

	t.Run("Can get posts list", func(t *testing.T) {
		actualPosts, err := posts.List(nil)
		if err != nil {
			t.Errorf("Failed to list posts %s", err)
		}

		if !calledPosts {
			t.Error("Didn't called posts")
		}

		if !reflect.DeepEqual(wantPosts, actualPosts) {
			t.Errorf("Posts aren't the same want: %+v, got: %+v", wantPosts, actualPosts)
		}
	})
}
