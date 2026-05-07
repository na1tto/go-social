package cache

import (
	"context"

	repository "github.com/na1tto/go-social/internal/store"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct{}

func (m MockUserStore) Get(ctx context.Context, id int64) (*repository.User, error) {
	return nil, nil
}

func (m MockUserStore) Set(ctx context.Context, user *repository.User) error {
	return nil
}
