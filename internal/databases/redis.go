package databases

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"test-task1/pkg/config"
)

func NewRedisClient(c *config.Cache) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", c.Host, c.Port),
		Password: c.Password,
		DB:       c.DBIndex,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
