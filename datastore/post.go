package datastore

import (
	"errors"

	"github.com/bogthe/bogthesrc"
)

func init() {
	DB.AddTableWithName(bogthesrc.Post{}, "post").SetKeys(false, "ID")
	createSql = append(createSql,
		`CREATE INDEX post_submittedat ON post(submittedat DESC);`,
	)
}

type PostStore struct {
	*Datastore
}

var (
	PostNotFound = errors.New("Post not found")
)

func (s *PostStore) Get(id string) (*bogthesrc.Post, error) {
	var posts []*bogthesrc.Post
	if err := s.dbh.Select(&posts, `SELECT * FROM post WHERE id=$1;`, id); err != nil {
		return nil, err
	}

	if len(posts) == 0 {
		return nil, PostNotFound
	}

	return posts[0], nil
}

func (s *PostStore) List(opt *bogthesrc.PostListOptions) ([]*bogthesrc.Post, error) {
	var posts []*bogthesrc.Post
	err := s.dbh.Select(&posts, `SELECT * FROM post LIMIT $1 OFFSET $2;`, opt.PerPageOrDefault(), opt.Offset())
	if err != nil {
		return nil, err
	}

	return posts, nil
}
