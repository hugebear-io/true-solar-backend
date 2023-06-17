package repo

import (
	"database/sql"

	"github.com/hugebear-io/true-solar-backend/internal/core/port"
)

type snmpConfigRepo struct {
	db *sql.DB
}

func NewSNMPConfigRepo(db *sql.DB) port.SNMPConfigRepoPort {
	return &snmpConfigRepo{db: db}
}

func (r snmpConfigRepo) GetSNMPConfig() (port.SNMPConfig, error) {
	queryString := `SELECT id, target_hostname, target_port, agent_hostname, created_at, updated_at
					FROM snmp_config
					LIMIT 1`
	row := r.db.QueryRow(queryString)
	if err := row.Err(); err != nil {
		return port.SNMPConfig{}, err
	}
	var result port.SNMPConfig
	err := row.Scan(
		&result.ID,
		&result.TargetHostname,
		&result.TargetPort,
		&result.AgentHostname,
		&result.CreatedAt,
		&result.UpdatedAt,
	)
	if err != nil {
		return port.SNMPConfig{}, err
	}
	return result, nil
}

func (r snmpConfigRepo) UpdateSNMPConfig(snmpConfig port.SNMPConfig) error {
	queryString := `UPDATE snmp_config
					SET target_hostname = ?,
						target_port = ?,
						agent_hostname = ?
					WHERE id = 1`
	stmt, err := r.db.Prepare(queryString)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		snmpConfig.TargetHostname,
		snmpConfig.TargetPort,
		snmpConfig.AgentHostname,
	)
	if err != nil {
		return err
	}
	return nil
}
