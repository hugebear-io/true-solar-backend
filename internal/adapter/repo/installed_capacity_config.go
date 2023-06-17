package repo

import (
	"database/sql"

	"github.com/hugebear-io/true-solar-backend/internal/core/port"
)

type installedCapacityConfigRepo struct {
	db *sql.DB
}

func NewInstalledCapacityConfigRepo(db *sql.DB) port.InstalledCapacityConfigRepoPort {
	return &installedCapacityConfigRepo{
		db: db,
	}
}

func (r installedCapacityConfigRepo) GetInstalledCapacityConfig() (port.InstalledCapacityConfig, error) {
	queryString := `SELECT id, efficiency_factor, focus_hour, created_at, updated_at
					FROM installed_capacity_config
					LIMIT 1`
	row := r.db.QueryRow(queryString)
	if err := row.Err(); err != nil {
		return port.InstalledCapacityConfig{}, err
	}
	var result port.InstalledCapacityConfig
	err := row.Scan(
		&result.ID,
		&result.EfficiencyFactor,
		&result.FocusHour,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return port.InstalledCapacityConfig{}, err
	}
	return result, nil
}

func (r installedCapacityConfigRepo) UpdateInstalledCapacityConfig(installedCapacityConfig port.InstalledCapacityConfig) error {
	queryString := `UPDATE installed_capacity_config
					SET efficiency_factor = ?,
						focus_hour = ?
					WHERE id = 1`
	stmt, err := r.db.Prepare(queryString)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		installedCapacityConfig.EfficiencyFactor,
		installedCapacityConfig.FocusHour,
	)
	if err != nil {
		return err
	}
	return nil
}
