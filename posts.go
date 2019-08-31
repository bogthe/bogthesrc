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

type PostListOptions struct {
	ListOptions
}

type PostService struct {
	Client *Client
}

func (s *PostService) List(options *PostListOptions) ([]*Post, error) {
	url, err := s.Client.url(router.Posts, nil, options)
	if err != nil {
		return nil, err
	}

	req, err := s.Client.NewRequest("GET", url.String())
	if err != nil {
		return nil, err
	}

	var posts []*Post
	_, err = s.Client.Do(req, &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (s *PostService) Get(id string) (*Post, error) {
	// need to have an url
	url, err := s.Client.url(router.Post, map[string]string{"ID": id}, nil)
	if err != nil {
		return nil, err
	}

	// create the request
	request, err := s.Client.NewRequest("GET", url.String())
	if err != nil {
		return nil, err
	}

	// client Do-es request with post as a body
	var post *Post
	_, err = s.Client.Do(request, &post)

	if err != nil {
		return nil, err
	}

	return post, nil
}
