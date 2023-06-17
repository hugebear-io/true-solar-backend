package domain

import "time"

type PerformanceAlarmConfigService interface {
	GetPerformanceAlarmConfig() (UpdatePerformanceAlarmRequestBody, error)
	UpdatePerformanceAlarmConfig(performanceAlarmConfig UpdatePerformanceAlarmRequestBody) error
}

type UpdatePerformanceAlarmRequestBody struct {
	Alarm1ID         int        `json:"alarm1_id"`
	Alarm1Name       string     `json:"alarm1_name"`
	Alarm1Interval   int        `json:"alarm1_interval" binding:"required"`
	Alarm1HitDay     int        `json:"alarm1_hit_day" binding:"required"`
	Alarm1Percentage float64    `json:"alarm1_percentage" binding:"required"`
	Alarm1Duration   int        `json:"alarm1_duration" binding:"required"`
	Alarm1CreatedAt  *time.Time `json:"alarm1_created_at"`
	Alarm1UpdatedAt  *time.Time `json:"alarm1_updated_at"`
	Alarm2ID         int        `json:"alarm2_id"`
	Alarm2Name       string     `json:"alarm2_name"`
	Alarm2Interval   int        `json:"alarm2_interval" binding:"required"`
	Alarm2Percentage float64    `json:"alarm2_percentage" binding:"required"`
	Alarm2Duration   int        `json:"alarm2_duration" binding:"required"`
	Alarm2CreatedAt  *time.Time `json:"alarm2_created_at"`
	Alarm2UpdatedAt  *time.Time `json:"alarm2_updated_at"`
}
