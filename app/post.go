package app

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/bogthe/bogthesrc"
	"github.com/bogthe/bogthesrc/router"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
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
	var opt bogthesrc.PostListOptions
	if err := schema.NewDecoder().Decode(&opt, r.URL.Query()); err != nil {
		return err
	}

	posts, err := apiClient.Posts.List(&opt)
	if err != nil {
		return err
	}

	return renderTemplate(w, r, PostListTemplate, http.StatusOK, struct {
		Posts []*bogthesrc.Post
	}{
		posts,
	})
}

func servePostCreate(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	schemaDecoder := schema.NewDecoder()

	var post bogthesrc.Post
	if err := schemaDecoder.Decode(&post, r.Form); err != nil {
		return err
	}

	if err := apiClient.Posts.Create(&post); err != nil {
		return err
	}

	redirectUrl := urlTo(router.Post, "ID", strconv.Itoa(post.ID))
	http.Redirect(w, r, redirectUrl.String(), http.StatusSeeOther)

	return nil
}

func servePostCreateForm(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query()

	post := &bogthesrc.Post{
		Title: getValueOrDefault(q, "Title"),
		Link:  getValueOrDefault(q, "Link"),
		Body:  getValueOrDefault(q, "Body"),
	}

	return renderTemplate(w, r, PostCreateTemplate, http.StatusOK, struct {
		Post         *bogthesrc.Post
		ActionTarget string
	}{
		post,
		router.PostCreate,
	})
}

func getValueOrDefault(q url.Values, name string) string {
	if val, success := q[name]; success {
		return val[0]
	}

	// url.Values is case-sensitive
	return q.Get(strings.ToLower(name))
}
