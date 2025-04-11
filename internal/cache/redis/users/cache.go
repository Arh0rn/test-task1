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

func New(client *redis.Client, ttl time.Duration) *UserCache {
	return &UserCache{
		client: client,
		ttl:    ttl,
	}
}

func (c *UserCache) Set(ctx context.Context, user *domain.User) error {
	key := userKey + fmt.Sprint(user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to marshal user", "error", err)
		return err
	}
	err = c.client.Set(ctx, key, data, c.ttl).Err()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to set user in cache", "error", err)
		return err
	}
	slog.DebugContext(ctx, "User set in cache", "user_id", user.ID)
	return nil
}

func (c *UserCache) GetAll(ctx context.Context) ([]*domain.User, error) {
	keys, err := c.client.Keys(ctx, "user:*").Result()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get keys from cache", "error", err)
		return nil, err
	}
	pipe := c.client.Pipeline()
	cmds := make([]*redis.StringCmd, len(keys))
	for i, key := range keys {
		cmds[i] = pipe.Get(ctx, key)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute pipeline", "error", err)
		return nil, err
	}

	var users []*domain.User
	for _, cmd := range cmds {
		data, err := cmd.Result()
		if err != nil {
			slog.ErrorContext(ctx, "Failed to get user from cache", "error", err)
			continue
		}
		var user domain.User
		if err := json.Unmarshal([]byte(data), &user); err == nil {
			users = append(users, &user)
		}
	}
	slog.DebugContext(ctx, "Users found in cache", "user_count", len(users))
	return users, nil
}

func (c *UserCache) SetAll(ctx context.Context, users []*domain.User) error {
	pipe := c.client.Pipeline()
	for _, user := range users {
		key := userKey + fmt.Sprint(user.ID)
		data, err := json.Marshal(user)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to marshal user", "error", err)
			continue
		}
		pipe.Set(ctx, key, data, c.ttl)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to execute pipeline", "error", err)
		return err
	}
	slog.DebugContext(ctx, "Users set in cache", "user_count", len(users))
	return nil
}

func (c *UserCache) GetByID(ctx context.Context, id int) (*domain.User, error) {
	key := userKey + fmt.Sprint(id)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			slog.InfoContext(ctx, "User not found in cache", "user_id", id)
			return nil, domain.ErrUserNotFound
		}
		slog.ErrorContext(ctx, "Failed to get user from cache", "error", err)
		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		slog.ErrorContext(ctx, "Failed to unmarshal user", "error", err)
		return nil, err

	}
	slog.DebugContext(ctx, "User found in cache", "user_id", user.ID)
	return &user, nil
}

func (c *UserCache) UpdateByID(ctx context.Context, update *domain.UserUpdate, id int) error {
	user, err := c.GetByID(ctx, id)
	if err != nil || user == nil {
		slog.ErrorContext(ctx, "Failed to get user for update", "error", err)
		return err
	}

	user.Name = update.Name
	user.Email = update.Email

	err = c.Set(ctx, user)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to set updated user in cache", "error", err)
		return err
	}

	slog.DebugContext(ctx, "User updated in cache", "user_id", user.ID)
	return nil
}

func (c *UserCache) DeleteByID(ctx context.Context, id int) error {
	key := userKey + fmt.Sprint(id)
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to delete user from cache", "error", err)
		return err
	}
	slog.DebugContext(ctx, "User deleted from cache", "user_id", id)
	return nil
}
