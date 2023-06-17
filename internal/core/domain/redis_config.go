package domain

import "github.com/hugebear-io/true-solar-backend/internal/core/port"

type RedisConfigService interface {
	GetRedisConfig() (port.RedisConfig, error)
	UpdateRedisConfig(redisConfig port.RedisConfig) error
}
