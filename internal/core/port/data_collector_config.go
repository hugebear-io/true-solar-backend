package port

type DataCollectorConfigRepoPort interface {
	GetAllDataCollectorConfig() ([]VendorAccount, error)
	GetOneDataCollectorConfig(id int) (VendorAccount, error)
	GetDataCollectorConfigByVendorType(vendorType string) ([]VendorAccount, error)
	CreateDataCollectorConfig(dataCollectorConfig VendorAccount) error
	UpdateDataCollectorConfig(id int, dataCollectorConfig VendorAccount) error
	DeleteDataCollectorConfig(id int) error
}
