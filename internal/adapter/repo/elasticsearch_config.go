package repo

import (
	"database/sql"

	"github.com/hugebear-io/true-solar-backend/internal/core/port"
)

type elasticSearchConfigRepo struct {
	db *sql.DB
}

func NewElasticSearchConfigRepo(db *sql.DB) port.ElasticSearchConfigRepoPort {
	return &elasticSearchConfigRepo{
		db: db,
	}
}

func (r elasticSearchConfigRepo) GetElasticsearchConfig() (port.ElasticSearchConfig, error) {
	queryString := `SELECT id, hostname, username, password, index_data, created_at, updated_at
					FROM elasticsearch_config
					LIMIT 1`
	row := r.db.QueryRow(queryString)
	if err := row.Err(); err != nil {
		return port.ElasticSearchConfig{}, err
	}
	var result port.ElasticSearchConfig
	err := row.Scan(
		&result.ID,
		&result.Hostname,
		&result.Username,
		&result.Password,
		&result.Index,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return port.ElasticSearchConfig{}, err
	}
	return result, nil
}

func (r elasticSearchConfigRepo) UpdateElasticsearchConfig(elasticsearchConfig port.ElasticSearchConfig) error {
	queryString := `UPDATE elasticsearch_config
					SET hostname = ?,
						username = ?,
						password = ?,
						index_data = ?
					WHERE id = 1`
	stmt, err := r.db.Prepare(queryString)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		elasticsearchConfig.Hostname,
		elasticsearchConfig.Username,
		elasticsearchConfig.Password,
		elasticsearchConfig.Index,
	)
	if err != nil {
		return err
	}
	return nil
}
