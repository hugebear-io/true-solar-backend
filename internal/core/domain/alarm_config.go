package domain

import "github.com/hugebear-io/true-solar-backend/internal/core/port"

type AlarmConfigService interface {
	GetAllAlarmConfig() ([]port.VendorAccount, error)
	GetAlarmConfigByVendorType(vendorType string) ([]port.VendorAccount, error)
	GetOneAlarmConfig(id int) (port.VendorAccount, error)
	CreateAlarmConfig(alarmConfig port.VendorAccount) error
	UpdateAlarmConfig(id int, alarmConfig port.VendorAccount) error
	DeleteAlarmConfig(id int) error
}
