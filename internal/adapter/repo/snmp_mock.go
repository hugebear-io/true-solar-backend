package repo

import "github.com/hugebear-io/true-solar-backend/internal/core/port"

type snmpRepoMock struct{}

func NewSNMPRepoMock() port.SNMPRepoPort {
	return &snmpRepoMock{}
}

func (r snmpRepoMock) SendAlarmTrap(deviceName string, alertName string, description string, severity string, lastedUpdateTime string) error {
	return nil
}

func (r snmpRepoMock) Close() {}
