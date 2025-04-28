package cache

import (
	"context"

	"example.com/goapi/internal/config"
	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

func NewClient(cfg *config.Conf) *Client {
	return &Client{
		redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr(),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		}),
	}
}

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.Client.Ping(ctx).Result()
	return err
}

func (c *Client) Close() error {
	if c.Client != nil {
		return c.Client.Close()
	}
	return nil
}
