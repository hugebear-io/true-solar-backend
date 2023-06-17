package infra

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/hugebear-io/true-solar-backend/pkg/config"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
)

func NewRedis(logger logger.Logger) *redis.Client {
	ctx := context.Background()
	cfg := config.Config.Redis
	options := redis.Options{
		Addr:     cfg.Addr,
		Username: cfg.User,
		Password: cfg.Password,
		DB:       cfg.DB,
	}
	client := redis.NewClient(&options)

	_, err := client.Ping(ctx).Result()
	if err != nil {
		logger.Panicf("NewRedis(): %v", err)
	}

	return client
}
