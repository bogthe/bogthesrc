package app

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/bogthe/bogthesrc"
	"github.com/bogthe/bogthesrc/router"
)

func TestPost(t *testing.T) {
	setup()
	defer teardown()

	called := false
	post := &bogthesrc.Post{ID: 1, Title: "Test title", Body: "Test body", Link: "testlink.com"}

	serviceMock := &bogthesrc.MockPostService{
		Get_: func(id int) (*bogthesrc.Post, error) {
			called = true
			if id != post.ID {
				t.Fatal("Post ID not found in test")
			}

			return post, nil
		},
	}

	apiClient = &bogthesrc.ApiClient{Posts: serviceMock}

	t.Run("Can display a single post", func(t *testing.T) {
		url, err := router.App().Get(router.Post).URL("ID", strconv.Itoa(post.ID))
		if err != nil {
			t.Fatal(err)
		}

		doc, resp := getHTML(t, url)

		if resp.Code != http.StatusOK {
			t.Errorf("Response code is wrong, got: %v", resp.Code)
		}

		if !called {
			t.Errorf("Handler not called")
		}

		a := doc.Find("a.post-link")
		if a.Text() != post.Title {
			t.Errorf("Post link text is wrong, wanted %s, got %s", post.Title, a.Text())
		}

		if got, _ := a.Attr("href"); got != post.Link {
			t.Errorf("Post link is wrong, wanted %s, got %s", post.Link, got)
		}

		body := doc.Find("p.post-body")
		if body.Text() != post.Body {
			t.Errorf("Body text is wrong, wanted %s, got %s", post.Body, body.Text())
		}
	})
}
