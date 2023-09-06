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
	addr, _ := redis.ParseURL(cfg.URI)
	client := redis.NewClient(addr)

	_, err := client.Ping(ctx).Result()
	if err != nil {
		logger.Panicf("NewRedis(): %v", err)
	}

	logger.Infof("NewRedis(): connected to %s", cfg.URI)
	return client
}
