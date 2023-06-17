package port

type AlarmConfigRepoPort interface {
	GetAllAlarmConfig() ([]VendorAccount, error)
	GetOneAlarmConfig(id int) (VendorAccount, error)
	GetAlarmConfigByVendorType(vendorType string) ([]VendorAccount, error)
	CreateAlarmConfig(alarmConfig VendorAccount) error
	UpdateAlarmConfig(id int, alarmConfig VendorAccount) error
	DeleteAlarmConfig(id int) error
}
