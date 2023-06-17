package port

import "time"

type PerformanceAlarmConfigRepoPort interface {
	GetPerformanceAlarmConfig() ([]PerformanceAlarmConfig, error)
	UpdatePerformanceAlarmConfig(id int, performanceAlarmConfig PerformanceAlarmConfig) error
}

type PerformanceAlarmConfig struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	Interval   int        `json:"interval"`
	HitDay     int        `json:"hit_day"`
	Percentage float64    `json:"percentage"`
	Duration   int        `json:"duration"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}
