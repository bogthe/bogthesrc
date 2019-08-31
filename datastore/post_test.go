package datastore

import (
	"reflect"
	"testing"

	"github.com/bogthe/bogthesrc"
)

func TestStoreGetWithDB(t *testing.T) {
	want := &bogthesrc.Post{ID: 1}

	tx, _ := DB.Begin()
	defer tx.Rollback()

	tx.Exec(`DELETE FROM post;`)
	if err := tx.Insert(want); err != nil {
		t.Fatalf("Failed insertion %s", err)
	}

	d := NewDatastore(tx)
	store := &PostStore{d}
	post, err := store.Get("1")
	if err != nil {
		t.Fatalf("Failed get %s", err)
	}

	if !reflect.DeepEqual(want, post) {
		t.Fatalf("Wanted %+v, got %+v", want, post)
	}
}

func TestStoreListWithDB(t *testing.T) {
	want := []*bogthesrc.Post{&bogthesrc.Post{ID: 1}}

	tx, _ := DB.Begin()
	defer tx.Rollback()

	tx.Exec(`DELETE FROM post;`)
	if err := tx.Insert(want[0]); err != nil {
		t.Fatalf("Failed insertion %s", err)
	}

	d := NewDatastore(tx)
	store := PostStore{d}
	posts, err := store.List(&bogthesrc.PostListOptions{ListOptions: bogthesrc.ListOptions{Page: 1, PerPage: 10}})
	if err != nil {
		t.Fatalf("Failed retrieve %s", err)
	}

	if !reflect.DeepEqual(want, posts) {
		t.Fatalf("Wanted %+v, got %+v", want, posts)
	}
}
