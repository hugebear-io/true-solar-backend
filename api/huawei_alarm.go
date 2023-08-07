package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/adapter/handler"
	"github.com/hugebear-io/true-solar-backend/internal/adapter/repo"
	"github.com/hugebear-io/true-solar-backend/internal/core/service"
	"github.com/hugebear-io/true-solar-backend/internal/infra"
	"github.com/hugebear-io/true-solar-backend/pkg/config"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
)

func BindHuaweiAlarmAPI(api *gin.RouterGroup) {
	cfg := config.Config
	l := logger.NewLogger(&logger.LoggerOption{
		LogName:     "logs/huawei-alarm.log",
		LogSize:     1024,
		LogAge:      90,
		LogBackup:   1,
		LogCompress: false,
		LogLevel:    logger.LogLevel(cfg.API.LogLevel),
		SkipCaller:  1,
	})

	rdb := infra.NewRedis(l)
	snmpClient := infra.NewSNMP(l)
	alarmConfigRepo := repo.NewAlarmConfigRepo(infra.SqlDB)
	alarmConfig := service.NewAlarmConfigService(alarmConfigRepo)

	snmp := repo.NewSNMPRepo(snmpClient, cfg.SNMP.AgentHost)
	serv := service.NewHuaweiAlarmService(alarmConfig, rdb, snmp, l)
	hdl := handler.NewHuaweiAlarmHandler(serv)

	sub := api.Group("/alarm")
	sub.GET("", hdl.Run)
}
