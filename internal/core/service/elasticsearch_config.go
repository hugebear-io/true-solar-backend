package service

import (
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
)

type elasticSearchConfigService struct {
	repo port.ElasticSearchConfigRepoPort
}

func NewElasticSearchConfigService(repo port.ElasticSearchConfigRepoPort) domain.ElasticSearchConfigService {
	return &elasticSearchConfigService{repo: repo}
}

func (s elasticSearchConfigService) GetElasticsearchConfig() (port.ElasticSearchConfig, error) {
	elasticsearchConfig, err := s.repo.GetElasticsearchConfig()
	if err != nil {
		return port.ElasticSearchConfig{}, err
	}
	return elasticsearchConfig, nil
}

func (s elasticSearchConfigService) UpdateElasticsearchConfig(elasticsearchConfig port.ElasticSearchConfig) error {
	err := s.repo.UpdateElasticsearchConfig(elasticsearchConfig)
	if err != nil {
		return err
	}
	return nil
}
