package port

import "time"

type SNMPConfigRepoPort interface {
	GetSNMPConfig() (SNMPConfig, error)
	UpdateSNMPConfig(snmpConfig SNMPConfig) error
}

type SNMPConfig struct {
	ID             int       `json:"id"`
	TargetHostname string    `json:"target_hostname" binding:"required"`
	TargetPort     string    `json:"target_port" binding:"required"`
	AgentHostname  string    `json:"agent_hostname" binding:"required"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
