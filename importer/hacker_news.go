package importer

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bogthe/bogthesrc"
)

type hackerNews struct {
	which string
}

const (
	base     = "https://hacker-news.firebaseio.com/v0/%s.json"
	describe = "https://hacker-news.firebaseio.com/v0/item/%v.json"
)

func init() {
	Fetchers = append(Fetchers, &hackerNews{"topstories"}, &hackerNews{"newstories"}, &hackerNews{"beststories"})
}

func (h *hackerNews) Fetch() ([]*bogthesrc.Post, error) {
	ids, err := getIds(fmt.Sprintf(base, h.which))
	if err != nil {
		return nil, err
	}

	posts := make([]*bogthesrc.Post, 0)
	for _, id := range ids {
		post, err := getPost(fmt.Sprintf(describe, id))
		if err != nil {
			return nil, err
		}

		if post == nil {
			continue
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (h *hackerNews) Site() string {
	return fmt.Sprintf("hackernews/%s", h.which)
}

func getPost(url string) (*bogthesrc.Post, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var hn *struct {
		Title string
		Score int
		Url   string
		Type  string
	}

	if err := json.NewDecoder(resp.Body).Decode(&hn); err != nil {
		return nil, err
	}

	if hn.Type != "story" || hn.Url == "" {
		return nil, nil
	}

	return &bogthesrc.Post{
		Title: hn.Title,
		Score: hn.Score,
		Link:  hn.Url,
	}, nil
}

func getIds(url string) ([]int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var ids []int
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, err
	}

	return ids, nil
}
