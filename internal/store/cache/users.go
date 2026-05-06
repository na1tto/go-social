package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	repository "github.com/na1tto/go-social/internal/store"
	"github.com/redis/go-redis/v9"
)

type UserStore struct {
	rdb *redis.Client
}

const UserExpTime = time.Minute * 3

func (s *UserStore) Get(ctx context.Context, userID int64) (*repository.User, error) {
	cacheKey := fmt.Sprintf("user=%v", userID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user repository.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (s *UserStore) Set(ctx context.Context, user *repository.User) error {
	if user.ID == 0 {
		return fmt.Errorf("user ID is required")
	}

	cacheKey := fmt.Sprintf("user=%v", user.ID)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.rdb.Set(ctx, cacheKey, json, UserExpTime).Err()
}
