package db

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"time"
	"wallet/internal/domain/wallet"
)

type Cache struct {
	logger *slog.Logger
	Client *redis.Client
}

func NewCache(logger *slog.Logger, c *redis.Client) wallet.Cache {
	return &Cache{
		logger: logger,
		Client: c,
	}
}

func (c *Cache) GetValue(ctx context.Context, key string) (string, error) {
	const op = "db.redis.GetValue"
	log := c.logger.With(slog.String("op", op))

	val, err := c.Client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		log.Error(err.Error())
		return "", err
	}
	return val, nil
}

func (c *Cache) SetValue(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	const op = "db.redis.SetValue"
	log := c.logger.With(slog.String("op", op))

	err := c.Client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
