package service

import (
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
)

type redisConfigService struct {
	repo port.RedisConfigRepoPort
}

func NewRedisConfigService(repo port.RedisConfigRepoPort) domain.RedisConfigService {
	return &redisConfigService{repo: repo}
}

func (s redisConfigService) GetRedisConfig() (port.RedisConfig, error) {
	redisConfig, err := s.repo.GetRedisConfig()
	if err != nil {
		return port.RedisConfig{}, err
	}
	return redisConfig, nil
}

func (s redisConfigService) UpdateRedisConfig(redisConfig port.RedisConfig) error {
	err := s.repo.UpdateRedisConfig(redisConfig)
	if err != nil {
		return err
	}
	return nil
}
