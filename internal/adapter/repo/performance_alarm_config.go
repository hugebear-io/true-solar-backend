package repo

import (
	"database/sql"

	"github.com/hugebear-io/true-solar-backend/internal/core/port"
)

type performanceAlarmConfigRepo struct {
	db *sql.DB
}

func NewPerformanceAlarmConfigRepo(db *sql.DB) port.PerformanceAlarmConfigRepoPort {
	return &performanceAlarmConfigRepo{
		db: db,
	}
}

func (r performanceAlarmConfigRepo) GetPerformanceAlarmConfig() ([]port.PerformanceAlarmConfig, error) {
	queryString := `SELECT id, name, interval, hit_day, percentage, duration, created_at, updated_at
					FROM performance_alarm_config
					ORDER BY id
					LIMIT 3`
	row, err := r.db.Query(queryString)
	if err != nil {
		return []port.PerformanceAlarmConfig{}, err
	}
	results := make([]port.PerformanceAlarmConfig, 0)
	for row.Next() {
		var result port.PerformanceAlarmConfig
		err := row.Scan(
			&result.ID,
			&result.Name,
			&result.Interval,
			&result.HitDay,
			&result.Percentage,
			&result.Duration,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return []port.PerformanceAlarmConfig{}, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (r performanceAlarmConfigRepo) UpdatePerformanceAlarmConfig(id int, performanceAlarmConfig port.PerformanceAlarmConfig) error {
	queryString := `UPDATE performance_alarm_config
					SET interval = ?,
						hit_day = ?,
						percentage = ?,
						duration = ?
					WHERE id = ?`
	stmt, err := r.db.Prepare(queryString)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		performanceAlarmConfig.Interval,
		performanceAlarmConfig.HitDay,
		performanceAlarmConfig.Percentage,
		performanceAlarmConfig.Duration,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}
