package bogthesrc

import (
	"time"

	"github.com/bogthe/bogthesrc/router"
)

type Post struct {
	ID          string
	Title       string
	Link        string
	Body        string
	SubmittedAt time.Time
	AuthordID   int
}

type PostService interface {
	Get(id string) (*Post, error)
}

type postService struct {
	client *Client
}

func (s *postService) Get(id string) (*Post, error) {
	// need to have an url
	url, err := s.client.url(router.Post, map[string]string{"ID": id})
	if err != nil {
		return nil, err
	}

	// create the request
	request, err := s.client.NewRequest("GET", url.String())
	if err != nil {
		return nil, err
	}

	// client Do-es request with post as a body
	var post *Post
	_, err = s.client.Do(request, &post)

	if err != nil {
		return nil, err
	}

	return post, nil
}
