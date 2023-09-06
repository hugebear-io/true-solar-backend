package service

import (
	"github.com/go-redis/redis/v8"
	"github.com/hugebear-io/true-solar-backend/internal/adapter/alarm"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
	"github.com/hugebear-io/true-solar/huawei"
)

type huaweiAlarmService struct {
	alarmConfig domain.AlarmConfigService
	rdb         *redis.Client
	snmp        port.SNMPRepoPort
	logger      logger.Logger
}

func NewHuaweiAlarmService(
	alarmConfig domain.AlarmConfigService,
	rdb *redis.Client,
	snmp port.SNMPRepoPort,
	logger logger.Logger,
) domain.HuaweiAlarmService {
	return &huaweiAlarmService{
		alarmConfig: alarmConfig,
		rdb:         rdb,
		snmp:        snmp,
		logger:      logger,
	}
}

func (s huaweiAlarmService) Run() error {
	usernames := make([]string, 0)
	password := ""

	configs, err := s.alarmConfig.GetAlarmConfigByVendorType(huawei.BRAND)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	for _, config := range configs {
		usernames = append(usernames, config.Username)
		password = config.Password
	}

	alarm := alarm.NewHuaweiAlarm(
		s.rdb,
		s.snmp,
		usernames,
		password,
		s.logger,
	)

	if err := alarm.Run(); err != nil {
		s.logger.Error(err)
		return err
	}

	return nil
}
