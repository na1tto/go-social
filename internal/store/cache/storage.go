package cache

import (
	"context"

	repository "github.com/na1tto/go-social/internal/store"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*repository.User, error)
		Set(context.Context, *repository.User) error
	}
}

func NewRedisStorage(rdb *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rdb: rdb},
	}
}
