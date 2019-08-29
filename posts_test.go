package bogthesrc

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/bogthe/bogthesrc/router"
)

func TestGETPost(t *testing.T) {
	setup()
	teardown()

	want := &Post{
		ID: "1",
	}

	called := false
	mux.HandleFunc(urlFor(t, router.Post, map[string]string{"ID": "1"}), func(w http.ResponseWriter, r *http.Request) {
		called = true

		checkMethod(t, r, "GET")
		writeJSON(w, want)
	})

	posts := &postService{client: client}
	actualPost, err := posts.Get("1")

	if err != nil {
		t.Errorf("Failed to get post :%s", err)
	}

	if !called {
		t.Error("Handler func not called")
	}

	if !reflect.DeepEqual(want, actualPost) {
		t.Errorf("Posts aren't the same want: %+v, got: %+v", want, actualPost)
	}
}
