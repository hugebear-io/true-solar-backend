package repo

import (
	"github.com/gosnmp/gosnmp"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
)

type snmpRepo struct {
	client    *gosnmp.GoSNMP
	agentHost string
	logger    logger.Logger
}

func NewSNMPRepo(client *gosnmp.GoSNMP, agentHost string) port.SNMPRepoPort {
	l := logger.NewLogger(&logger.LoggerOption{
		LogName:     "logs/snmp.log",
		LogSize:     1024,
		LogAge:      90,
		LogBackup:   1,
		LogCompress: false,
		LogLevel:    logger.LogLevel(logger.LOG_LEVEL_DEBUG),
		SkipCaller:  1,
	})

	return &snmpRepo{client: client, agentHost: agentHost, logger: l}
}

func (r snmpRepo) SendAlarmTrap(deviceName string, alertName string, description string, severity string, lastedUpdateTime string) error {
	pduClass := gosnmp.SnmpPDU{
		Name:  "1.3.6.1.4.1.30378.2.1",
		Type:  gosnmp.OctetString,
		Value: "HPOVComponent",
	}
	pduName := gosnmp.SnmpPDU{
		Name:  "1.3.6.1.4.1.30378.2.2",
		Type:  gosnmp.OctetString,
		Value: deviceName,
	}
	pduAlert := gosnmp.SnmpPDU{
		Name:  "1.3.6.1.4.1.30378.2.3",
		Type:  gosnmp.OctetString,
		Value: alertName,
	}
	pduDesc := gosnmp.SnmpPDU{
		Name:  "1.3.6.1.4.1.30378.2.4",
		Type:  gosnmp.OctetString,
		Value: description,
	}
	pduSeverity := gosnmp.SnmpPDU{
		Name:  "1.3.6.1.4.1.30378.2.5",
		Type:  gosnmp.OctetString,
		Value: severity,
	}
	pduLastedUpdateTime := gosnmp.SnmpPDU{
		Name:  "1.3.6.1.4.1.30378.2.6",
		Type:  gosnmp.OctetString,
		Value: lastedUpdateTime,
	}
	trap := gosnmp.SnmpTrap{
		Enterprise:   "1.3.6.1.4.1.30378.1.1",
		AgentAddress: r.agentHost,
		GenericTrap:  6,
		SpecificTrap: 1,
		Variables:    []gosnmp.SnmpPDU{pduClass, pduName, pduAlert, pduDesc, pduSeverity, pduLastedUpdateTime},
	}

	result, err := r.client.SendTrap(trap)
	if err != nil {
		r.logger.Error(err)
		return err
	}

	r.logger.Infof("Trap: %#v", trap)
	r.logger.Infof("Result: %#v", result)
	return nil
}

func (r *snmpRepo) Close() {
	r.client.Conn.Close()
}
