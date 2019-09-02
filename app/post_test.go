package app

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/bogthe/bogthesrc"
	"github.com/bogthe/bogthesrc/router"
)

func TestPost(t *testing.T) {
	setup()
	defer teardown()

	called := false
	listCalled := false
	createCalled := false
	post := &bogthesrc.Post{ID: 1, Title: "Test title", Body: "Test body", Link: "testlink.com"}
	posts := []*bogthesrc.Post{post}

	serviceMock := &bogthesrc.MockPostService{
		Get_: func(id int) (*bogthesrc.Post, error) {
			called = true
			if id != post.ID {
				t.Fatal("Post ID not found in test")
			}

			return post, nil
		},
		Create_: func(post *bogthesrc.Post) error {
			createCalled = true
			post.ID = 1
			return nil
		},
		List_: func(opt *bogthesrc.PostListOptions) ([]*bogthesrc.Post, error) {
			listCalled = true
			return posts, nil
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

		checkPostRender(doc, post, t)
	})

	t.Run("Can display a list of posts", func(t *testing.T) {
		url, err := router.App().Get(router.Posts).URL()
		if err != nil {
			t.Fatal(err)
		}

		doc, resp := getHTML(t, url)

		if resp.Code != http.StatusOK {
			t.Errorf("Response code is wrong, got: %v", resp.Code)
		}

		if !listCalled {
			t.Errorf("Handler not called")
		}

		for _, p := range posts {
			checkPostRender(doc, p, t)
		}
	})

	t.Run("Can get post creation form", func(t *testing.T) {
		url_, err := router.App().Get(router.PostCreateForm).URL()
		if err != nil {
			t.Fatal(err)
		}

		url_.RawQuery = url.Values{
			"Title": []string{post.Title},
			"Link":  []string{post.Link},
			"Body":  []string{post.Body},
		}.Encode()

		doc, resp := getHTML(t, url_)
		if resp.Code != http.StatusOK {
			t.Fatal("Response code NOT OK")
		}

		if got, _ := doc.Find("input[name=Title]").Attr("value"); got != post.Title {
			t.Errorf("Form input expected %s, got %s", post.Title, got)
		}

		if got, _ := doc.Find("input[name=Link]").Attr("value"); got != post.Link {
			t.Errorf("Form input expected %s, got %s", post.Link, got)
		}

		if got := doc.Find("textarea[name=Body]").Text(); got != post.Body {
			t.Errorf("Form input expected %s, got %s", post.Body, got)
		}
	})

	t.Run("Can redirect after creating a post", func(t *testing.T) {
		url_, err := router.App().Get(router.Posts).URL()
		if err != nil {
			t.Fatal(err)
		}

		v := url.Values{
			"Title": []string{post.Title},
			"Link":  []string{post.Link},
			"Body":  []string{post.Body},
		}

		req, err := http.NewRequest("POST", url_.String(), strings.NewReader(v.Encode()))
		if err != nil {
			t.Fatal(err)
		}

		resp := httptest.NewRecorder()
		resp.Body = new(bytes.Buffer)
		testMux.ServeHTTP(resp, req)

		if !createCalled {
			t.Error("Client method not called")
		}

		if resp.Code != http.StatusSeeOther {
			t.Errorf("Wrong response code %v", resp.Code)
		}

		redirectUrl := urlTo(router.Post, "ID", strconv.Itoa(post.ID)).String()
		if resp.Header().Get("location") != redirectUrl {
			t.Errorf("Bad location %v", resp.Header().Get("location"))
		}
	})
}

func checkFormValue(input, expected string, doc *goquery.Document, t *testing.T) {
	if got, _ := doc.Find(input).Attr("value"); got != expected {
		t.Errorf("Form input expected %s, got %s", expected, input)
	}
}

func checkPostRender(doc *goquery.Document, post *bogthesrc.Post, t *testing.T) {
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
}
