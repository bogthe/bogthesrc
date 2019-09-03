package datastore

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bogthe/bogthesrc"
	"github.com/jmoiron/modl"
)

func init() {
	DB.AddTableWithName(bogthesrc.Post{}, "post").SetKeys(true, "ID")
	createSql = append(createSql,
		`CREATE INDEX post_submittedat ON post(submittedat DESC);`,
		`CREATE UNIQUE INDEX post_linkurl ON post(link);`,
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

func (s *PostStore) Create(post *bogthesrc.Post) error {
	retry := 0
	shouldRetry := false
	var reason error

	for retry < 3 {
		err := transact(s.dbh, func(tx modl.SqlExecutor) error {
			var existing []*bogthesrc.Post
			if err := tx.Select(&existing, `SELECT * FROM post WHERE link=$1 LIMIT 1;`, post.Link); err != nil {
				return err
			}

			if len(existing) > 0 {
				*post = *existing[0]
				return nil
			}

			if err := tx.Insert(post); err != nil {
				if strings.Contains(err.Error(), `violates unique constraint "post_linkurl"`) {
					time.Sleep(time.Duration(rand.Intn(75)) * time.Millisecond)
					shouldRetry = true
					return err
				}
				return err
			}

			return nil
		})

		if !shouldRetry {
			return nil
		}

		if err != nil && reason == nil {
			reason = err
		}

		retry++
	}

	return fmt.Errorf("Couldn't insert post %+v, reason: %s", post, reason)
}

func (s *PostStore) List(opt *bogthesrc.PostListOptions) ([]*bogthesrc.Post, error) {
	var posts []*bogthesrc.Post
	err := s.dbh.Select(&posts, `SELECT * FROM post LIMIT $1 OFFSET $2;`, opt.PerPageOrDefault(), opt.Offset())
	if err != nil {
		return nil, err
	}

	return posts, nil
}
