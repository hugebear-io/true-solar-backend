package domain

import "github.com/hugebear-io/true-solar-backend/internal/core/port"

type DataCollectorConfigService interface {
	GetAllDataCollectorConfig() ([]port.VendorAccount, error)
	GetDataCollectorConfigByVendorType(vendorType string) ([]port.VendorAccount, error)
	GetOneDataCollectorConfig(id int) (port.VendorAccount, error)
	CreateDataCollectorConfig(dataCollectorConfig port.VendorAccount) error
	UpdateDataCollectorConfig(id int, dataCollectorConfig port.VendorAccount) error
	DeleteDataCollectorConfig(id int) error
}
