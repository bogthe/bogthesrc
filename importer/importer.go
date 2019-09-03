package importer

import (
	"github.com/bogthe/bogthesrc"
)

type Fetcher interface {
	Fetch() ([]*bogthesrc.Post, error)
	Site() string
}

var (
	Fetchers []Fetcher

	apiClient = bogthesrc.NewApiClient(nil)
)

func Import(f Fetcher) []error {
	posts, err := f.Fetch()
	if err != nil {
		return []error{err}
	}

	errors := make([]error, 0)
	for _, post := range posts {
		err := apiClient.Posts.Create(post)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}
