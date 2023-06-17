package service

import (
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
)

type installedCapacityConfigService struct {
	repo port.InstalledCapacityConfigRepoPort
}

func NewInstalledCapacityConfigService(repo port.InstalledCapacityConfigRepoPort) domain.InstalledCapacityConfigService {
	return &installedCapacityConfigService{repo: repo}
}

func (s installedCapacityConfigService) GetInstalledCapacityConfig() (port.InstalledCapacityConfig, error) {
	installedCapacityConfig, err := s.repo.GetInstalledCapacityConfig()
	if err != nil {
		return port.InstalledCapacityConfig{}, err
	}
	return installedCapacityConfig, nil
}

func (s installedCapacityConfigService) UpdateInstalledCapacityConfig(installedCapacityConfig port.InstalledCapacityConfig) error {
	err := s.repo.UpdateInstalledCapacityConfig(installedCapacityConfig)
	if err != nil {
		return err
	}
	return nil
}
