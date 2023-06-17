package domain

import "github.com/hugebear-io/true-solar-backend/internal/core/port"

type InstalledCapacityConfigService interface {
	GetInstalledCapacityConfig() (port.InstalledCapacityConfig, error)
	UpdateInstalledCapacityConfig(installedCapacityConfig port.InstalledCapacityConfig) error
}
