package repo

import (
	"database/sql"

	"github.com/hugebear-io/true-solar-backend/internal/core/port"
)

type dataCollectorConfigRepo struct {
	db *sql.DB
}

func NewDataCollectorConfigRepo(db *sql.DB) port.DataCollectorConfigRepoPort {
	return &dataCollectorConfigRepo{
		db: db,
	}
}

func (r dataCollectorConfigRepo) GetAllDataCollectorConfig() ([]port.VendorAccount, error) {
	queryString := `SELECT id, vendor_type, username, password, app_id, app_secret, token, created_at, updated_at
					FROM data_collector_config`
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

func (r dataCollectorConfigRepo) GetOneDataCollectorConfig(id int) (port.VendorAccount, error) {
	queryString := `SELECT id, vendor_type, username, password, app_id, app_secret, token, created_at, updated_at
					FROM data_collector_config
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

func (r dataCollectorConfigRepo) GetDataCollectorConfigByVendorType(vendorType string) ([]port.VendorAccount, error) {
	queryString := `SELECT id, vendor_type, username, password, app_id, app_secret, token, created_at, updated_at
					FROM data_collector_config
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

func (r dataCollectorConfigRepo) CreateDataCollectorConfig(dataCollectorConfig port.VendorAccount) error {
	queryString := `INSERT INTO data_collector_config (vendor_type, username, password, app_id, app_secret, token)
					VALUES (?, ?, ?, ?, ?, ?)`
	stmt, err := r.db.Prepare(queryString)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		dataCollectorConfig.VendorType,
		dataCollectorConfig.Username,
		dataCollectorConfig.Password,
		dataCollectorConfig.AppID,
		dataCollectorConfig.AppSecret,
		dataCollectorConfig.Token,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r dataCollectorConfigRepo) UpdateDataCollectorConfig(id int, dataCollectorConfig port.VendorAccount) error {
	queryString := `UPDATE data_collector_config
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
		dataCollectorConfig.VendorType,
		dataCollectorConfig.Username,
		dataCollectorConfig.Password,
		dataCollectorConfig.AppID,
		dataCollectorConfig.AppSecret,
		dataCollectorConfig.Token,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r dataCollectorConfigRepo) DeleteDataCollectorConfig(id int) error {
	queryString := `DELETE FROM data_collector_config
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
