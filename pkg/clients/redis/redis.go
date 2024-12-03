package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type ConfigRedis struct {
	Addr     string
	Username string
	Password string
	DB       string
}

func NewClient(cfg ConfigRedis) (*redis.Client, error) {
	db, err := strconv.Atoi(cfg.DB)
	if err != nil {
		return nil, err
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       db,
	})
	code := rdb.Ping(context.Background())
	if err := code.Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}
