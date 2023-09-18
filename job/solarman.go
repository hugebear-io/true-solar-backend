package job

import (
	"github.com/hugebear-io/true-solar-backend/internal/adapter/handler"
	"github.com/hugebear-io/true-solar-backend/internal/adapter/repo"
	"github.com/hugebear-io/true-solar-backend/internal/core/service"
	"github.com/hugebear-io/true-solar-backend/internal/infra"
)

func NewSolarmanJobHandler() (*handler.SolarmanCollectorHandler, *handler.SolarmanAlarmHandler) {
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

	return collectorHdl, alarmHdl
}
