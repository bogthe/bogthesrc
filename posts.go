package bogthesrc

import (
	"strconv"
	"time"

	"github.com/bogthe/bogthesrc/router"
)

type Post struct {
	ID          int
	Title       string
	Link        string
	Body        string
	SubmittedAt time.Time
	AuthordID   int
}

type PostListOptions struct {
	ListOptions
}

type PostService interface {
	List(options *PostListOptions) ([]*Post, error)
	Create(post *Post) error
	Get(id int) (*Post, error)
}

type postService struct {
	Client *Client
}

func (s *postService) List(options *PostListOptions) ([]*Post, error) {
	url, err := s.Client.url(router.Posts, nil, options)
	if err != nil {
		return nil, err
	}

	req, err := s.Client.NewRequest("GET", url.String(), nil)
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

func (s *postService) Create(post *Post) error {
	url, err := s.Client.url(router.PostCreate, nil, nil)
	if err != nil {
		return err
	}

	req, err := s.Client.NewRequest("POST", url.String(), post)
	if err != nil {
		return err
	}

	_, err = s.Client.Do(req, &post)
	if err != nil {
		return err
	}

	return nil
}

func (s *postService) Get(id int) (*Post, error) {
	// need to have an url
	url, err := s.Client.url(router.Post, map[string]string{"ID": strconv.Itoa(id)}, nil)
	if err != nil {
		return nil, err
	}

	// create the request
	request, err := s.Client.NewRequest("GET", url.String(), nil)
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
