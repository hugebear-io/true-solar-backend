package service

import (
	"github.com/hugebear-io/true-solar-backend/internal/adapter/alarm"
	"github.com/hugebear-io/true-solar-backend/internal/adapter/repo"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/infra"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
)

type solarmanAlarmService struct {
	alarmConfig domain.AlarmConfigService
	logger      logger.Logger
}

func NewSolarmanAlarmService(alarmConfig domain.AlarmConfigService) domain.SolarmanAlarmService {
	l := logger.NewLogger(&logger.LoggerOption{
		LogName:     "logs/solarman-alarm.log",
		LogSize:     1024,
		LogAge:      90,
		LogBackup:   1,
		LogCompress: false,
		LogLevel:    logger.LogLevel(logger.LOG_LEVEL_DEBUG),
		SkipCaller:  1,
	})

	return &solarmanAlarmService{
		alarmConfig: alarmConfig,
		logger:      l,
	}
}

func (s solarmanAlarmService) Run() error {
	rdb := infra.NewRedis(s.logger)
	snmp := infra.NewSNMP(s.logger)

	defer rdb.Close()
	defer snmp.Conn.Close()

	configs, err := s.alarmConfig.GetAlarmConfigByVendorType(constant.VENDOR_TYPE_INVT)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	usernames := []string{}
	password := ""
	appID := ""
	appSecret := ""
	for _, config := range configs {
		usernames = append(usernames, config.Username)
		password = config.Password
		appID = *config.AppID
		appSecret = *config.AppSecret
	}

	usernames = []string{"bignode.invt.th@gmail.com"}
	password = "123456*"
	appID = "202010143565002"
	appSecret = "222c202135013aee622c71cdf8c47757"

	// snmpRepo := repo.NewSNMPRepo(snmp, "")
	snmpRepo := repo.NewSNMPRepoMock()
	solarmanAlarm := alarm.NewSolarmanAlarm(rdb, snmpRepo, usernames, password, appID, appSecret)
	solarmanAlarm.Run()

	return nil
}

func (s *solarmanAlarmService) Close() {
	s.logger.Close()
}
