package port

import "time"

type ElasticSearchConfigRepoPort interface {
	GetElasticsearchConfig() (ElasticSearchConfig, error)
	UpdateElasticsearchConfig(elasticsearchConfig ElasticSearchConfig) error
}

type ElasticSearchConfig struct {
	ID        int       `json:"id"`
	Hostname  string    `json:"hostname" binding:"required"`
	Username  string    `json:"username" binding:"required"`
	Password  string    `json:"password" binding:"required"`
	Index     string    `json:"index" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
