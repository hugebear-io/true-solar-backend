package port

import "time"

type InstalledCapacityConfigRepoPort interface {
	GetInstalledCapacityConfig() (InstalledCapacityConfig, error)
	UpdateInstalledCapacityConfig(installedCapacityConfig InstalledCapacityConfig) error
}

type InstalledCapacityConfig struct {
	ID               int       `json:"id"`
	EfficiencyFactor float64   `json:"efficiency_factor" binding:"required"`
	FocusHour        int       `json:"focus_hour" binding:"required"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
