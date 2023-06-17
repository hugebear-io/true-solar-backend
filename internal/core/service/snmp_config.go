package service

import (
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
)

type snmpConfigService struct {
	repo port.SNMPConfigRepoPort
}

func NewSNMPConfigService(repo port.SNMPConfigRepoPort) domain.SNMPConfigService {
	return &snmpConfigService{repo: repo}
}

func (s snmpConfigService) GetSNMPConfig() (port.SNMPConfig, error) {
	snmpConfig, err := s.repo.GetSNMPConfig()
	if err != nil {
		return port.SNMPConfig{}, err
	}
	return snmpConfig, nil
}

func (s snmpConfigService) UpdateSNMPConfig(snmpConfig port.SNMPConfig) error {
	err := s.repo.UpdateSNMPConfig(snmpConfig)
	if err != nil {
		return err
	}
	return nil
}
