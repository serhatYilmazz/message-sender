package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/serhatYilmazz/message-sender/internal/config"
	"github.com/sirupsen/logrus"
)

type Client struct {
	*redis.Client
	logger *logrus.Logger
}

func NewClient(cfg config.RedisConfig, logger *logrus.Logger) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolSize:     cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.WithError(err).Error("failed to connect to Redis")
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("successfully connected to Redis")

	return &Client{
		Client: rdb,
		logger: logger,
	}, nil
}

func (c *Client) Close() error {
	c.logger.Info("closing Redis connection")
	return c.Client.Close()
}
