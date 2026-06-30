package cache

import (
	"context"
	"delivery-tracker/internal/config"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func NewClient(cfg config.RedisConfig) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %v", err)
	}

	return rdb, nil
}
