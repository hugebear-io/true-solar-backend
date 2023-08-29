package infra

import (
	"strconv"
	"time"

	"github.com/gosnmp/gosnmp"
	"github.com/hugebear-io/true-solar-backend/pkg/config"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
)

func NewSNMP(logger logger.Logger) *gosnmp.GoSNMP {
	cfg := config.Config.SNMP
	port, err := strconv.Atoi(cfg.TargetPort)
	if err != nil {
		logger.Panic(err)
	}

	SNMP := &gosnmp.GoSNMP{
		Port:               uint16(port),
		Transport:          "udp",
		Community:          "public",
		Version:            gosnmp.Version1,
		Timeout:            time.Duration(300) * time.Second,
		Retries:            20,
		ExponentialTimeout: true,
		MaxOids:            gosnmp.MaxOids,
		Target:             cfg.TargetHost,
	}

	if err := SNMP.Connect(); err != nil {
		logger.Panic(err)
		return nil
	}

	logger.Info("Initialized SNMP")
	return SNMP
}
