package domain

import "github.com/hugebear-io/true-solar-backend/internal/core/port"

type ElasticSearchConfigService interface {
	GetElasticsearchConfig() (port.ElasticSearchConfig, error)
	UpdateElasticsearchConfig(elasticsearchConfig port.ElasticSearchConfig) error
}
