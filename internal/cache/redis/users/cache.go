package users

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"test-task1/internal/domain"
	"time"
)

const (
	userKey = "user:"
)

type UserCache struct {
	client *redis.Client
	ttl    time.Duration
}

func (c *UserCache) Set(ctx context.Context, user *domain.User) error {
	key := userKey + fmt.Sprint(user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to marshal user", "error", err)
		return err
	}
	return c.client.Set(ctx, key, data, c.ttl).Err()
}

func (c *UserCache) GetAll(ctx context.Context) ([]*domain.User, error) {
	keys, err := c.client.Keys(ctx, "user:*").Result()
	if err != nil {
		return nil, err
	}

	pipe := c.client.Pipeline()
	cmds := make([]*redis.StringCmd, len(keys))
	for i, key := range keys {
		cmds[i] = pipe.Get(ctx, key)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	var users []*domain.User
	for _, cmd := range cmds {
		data, err := cmd.Result()
		if err != nil {
			continue
		}
		var user domain.User
		if err := json.Unmarshal([]byte(data), &user); err == nil {
			users = append(users, &user)
		}
	}
	return users, nil
}

func (c *UserCache) SetAll(ctx context.Context, users []*domain.User) error {
	pipe := c.client.Pipeline()
	for _, user := range users {
		key := userKey + fmt.Sprint(user.ID)
		data, err := json.Marshal(user)
		if err != nil {
			continue
		}
		pipe.Set(ctx, key, data, c.ttl)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (c *UserCache) GetByID(ctx context.Context, id int) (*domain.User, error) {
	key := userKey + fmt.Sprint(id)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			slog.InfoContext(ctx, "User not found in cache", "user_id", id)
			return nil, nil
		}
		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		slog.ErrorContext(ctx, "Failed to unmarshal user", "error", err)
		return nil, err

	}
	return &user, nil
}

func (c *UserCache) UpdateByID(ctx context.Context, update *domain.UserUpdate, id int) error {
	user, err := c.GetByID(ctx, id)
	if err != nil || user == nil {
		return err
	}

	user.Name = update.Name
	user.Email = update.Email

	return c.Set(ctx, user)
}

func (c *UserCache) DeleteByID(ctx context.Context, id int) error {
	key := userKey + fmt.Sprint(id)
	return c.client.Del(ctx, key).Err()
}

func New(client *redis.Client, ttl time.Duration) *UserCache {
	return &UserCache{
		client: client,
		ttl:    ttl,
	}
}
