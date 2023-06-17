package domain

import "github.com/hugebear-io/true-solar-backend/internal/core/port"

type SNMPConfigService interface {
	GetSNMPConfig() (port.SNMPConfig, error)
	UpdateSNMPConfig(snmpConfig port.SNMPConfig) error
}
