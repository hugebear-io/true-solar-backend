package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/adapter/handler"
	"github.com/hugebear-io/true-solar-backend/internal/adapter/repo"
	"github.com/hugebear-io/true-solar-backend/internal/core/service"
	"github.com/hugebear-io/true-solar-backend/internal/infra"
)

func BindSolarmanAPI(api *gin.RouterGroup) {
	siteRegionMappingRepo := repo.NewSiteRegionMappingRepo(infra.SqlDB)
	siteRegionMapping := service.NewSiteRegionMappingService(siteRegionMappingRepo)

	dataCollectorConfigRepo := repo.NewDataCollectorConfigRepo(infra.SqlDB)
	dataCollectorConfig := service.NewDataCollectorConfigService(dataCollectorConfigRepo)

	alarmConfigRepo := repo.NewAlarmConfigRepo(infra.SqlDB)
	alarmConfig := service.NewAlarmConfigService(alarmConfigRepo)

	collectorServ := service.NewSolarmanCollectorService(
		dataCollectorConfig,
		siteRegionMapping,
	)
	collectorHdl := handler.NewSolarmanCollectorHandler(collectorServ)

	alarmServ := service.NewSolarmanAlarmService(alarmConfig)
	alarmHdl := handler.NewSolarmanAlarmHandler(alarmServ)

	sub := api.Group("")
	sub.GET("/collector", collectorHdl.Run)
	sub.GET("/alarm", alarmHdl.Run)
}
