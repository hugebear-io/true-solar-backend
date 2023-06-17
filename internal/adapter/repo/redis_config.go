package repo

import (
	"database/sql"

	"github.com/hugebear-io/true-solar-backend/internal/core/port"
)

type redisConfigRepo struct {
	db *sql.DB
}

func NewRedisConfigRepo(db *sql.DB) port.RedisConfigRepoPort {
	return &redisConfigRepo{
		db: db,
	}
}
func (r redisConfigRepo) GetRedisConfig() (port.RedisConfig, error) {
	queryString := `SELECT id, hostname, username, password, created_at, updated_at
					FROM redis_config
					LIMIT 1`
	row := r.db.QueryRow(queryString)
	if err := row.Err(); err != nil {
		return port.RedisConfig{}, err
	}
	var result port.RedisConfig
	err := row.Scan(
		&result.ID,
		&result.Hostname,
		&result.Username,
		&result.Password,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return port.RedisConfig{}, err
	}
	return result, nil
}

func (r redisConfigRepo) UpdateRedisConfig(redisConfig port.RedisConfig) error {
	queryString := `UPDATE redis_config
					SET hostname = ?,
						username = ?,
						password = ?
					WHERE id = 1`
	stmt, err := r.db.Prepare(queryString)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		redisConfig.Hostname,
		redisConfig.Username,
		redisConfig.Password,
	)
	if err != nil {
		return err
	}
	return nil
}
