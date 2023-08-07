package port

type SNMPRepoPort interface {
	SendAlarmTrap(deviceName string, alertName string, description string, severity string, lastedUpdateTime string) error
	Close()
}
