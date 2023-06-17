package repo

import (
	"database/sql"

	"github.com/hugebear-io/true-solar-backend/internal/core/port"
)

type alarmConfigRepo struct {
	db *sql.DB
}

func NewAlarmConfigRepo(db *sql.DB) port.AlarmConfigRepoPort {
	return &alarmConfigRepo{
		db: db,
	}
}

func (r alarmConfigRepo) GetAllAlarmConfig() ([]port.VendorAccount, error) {
	queryString := `SELECT id, vendor_type, username, password, app_id, app_secret, token, created_at, updated_at
					FROM alarm_config`
	row, err := r.db.Query(queryString)
	if err != nil {
		return []port.VendorAccount{}, err
	}
	results := make([]port.VendorAccount, 0)
	for row.Next() {
		var result port.VendorAccount
		err := row.Scan(
			&result.ID,
			&result.VendorType,
			&result.Username,
			&result.Password,
			&result.AppID,
			&result.AppSecret,
			&result.Token,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return []port.VendorAccount{}, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (r alarmConfigRepo) GetOneAlarmConfig(id int) (port.VendorAccount, error) {
	queryString := `SELECT id, vendor_type, username, password, app_id, app_secret, token, created_at, updated_at
					FROM alarm_config
					WHERE id = ?`
	row := r.db.QueryRow(queryString, id)
	if err := row.Err(); err != nil {
		return port.VendorAccount{}, err
	}
	var result port.VendorAccount
	err := row.Scan(
		&result.ID,
		&result.VendorType,
		&result.Username,
		&result.Password,
		&result.AppID,
		&result.AppSecret,
		&result.Token,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return port.VendorAccount{}, err
	}
	return result, nil
}

func (r alarmConfigRepo) GetAlarmConfigByVendorType(vendorType string) ([]port.VendorAccount, error) {
	queryString := `SELECT id, vendor_type, username, password, app_id, app_secret, token, created_at, updated_at
					FROM alarm_config
					WHERE vendor_type = ?`
	row, err := r.db.Query(queryString, vendorType)
	if err != nil {
		return []port.VendorAccount{}, err
	}
	results := make([]port.VendorAccount, 0)
	for row.Next() {
		var result port.VendorAccount
		err := row.Scan(
			&result.ID,
			&result.VendorType,
			&result.Username,
			&result.Password,
			&result.AppID,
			&result.AppSecret,
			&result.Token,
			&result.CreatedAt,
			&result.UpdatedAt,
		)
		if err != nil {
			return []port.VendorAccount{}, err
		}
		results = append(results, result)
	}
	return results, nil
}

func (r alarmConfigRepo) CreateAlarmConfig(alarmConfig port.VendorAccount) error {
	queryString := `INSERT INTO alarm_config (vendor_type, username, password, app_id, app_secret, token)
					VALUES (?, ?, ?, ?, ?, ?)`
	stmt, err := r.db.Prepare(queryString)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		alarmConfig.VendorType,
		alarmConfig.Username,
		alarmConfig.Password,
		alarmConfig.AppID,
		alarmConfig.AppSecret,
		alarmConfig.Token,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r alarmConfigRepo) UpdateAlarmConfig(id int, alarmConfig port.VendorAccount) error {
	queryString := `UPDATE alarm_config
					SET vendor_type = ?,
						username = ?,
						password = ?,
						app_id = ?,
						app_secret = ?,
						token = ?
					WHERE id = ?`
	stmt, err := r.db.Prepare(queryString)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		alarmConfig.VendorType,
		alarmConfig.Username,
		alarmConfig.Password,
		alarmConfig.AppID,
		alarmConfig.AppSecret,
		alarmConfig.Token,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r alarmConfigRepo) DeleteAlarmConfig(id int) error {
	queryString := `DELETE FROM alarm_config
					WHERE id = ?`
	stmt, err := r.db.Prepare(queryString)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
