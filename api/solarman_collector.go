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

func BindSolarmanCollectorAPI(api *gin.RouterGroup) {
	apiConfig := config.Config.API
	elasticConfig := config.Config.ElasticSearch

	l := logger.NewLogger(&logger.LoggerOption{
		LogName:     "logs/solarman-inverter.log",
		LogSize:     1024,
		LogAge:      90,
		LogBackup:   1,
		LogCompress: false,
		LogLevel:    logger.LogLevel(apiConfig.LogLevel),
		SkipCaller:  1,
	})

	elastic := repo.NewElasticSearchRepo(infra.ElasticSearch, elasticConfig.Index)

	siteRegionMappingRepo := repo.NewSiteRegionMappingRepo(infra.SqlDB)
	siteRegionMapping := service.NewSiteRegionMappingService(siteRegionMappingRepo)

	dataCollectorConfigRepo := repo.NewDataCollectorConfigRepo(infra.SqlDB)
	dataCollectorConfig := service.NewDataCollectorConfigService(dataCollectorConfigRepo)

	serv := service.NewSolarmanCollectorService(
		dataCollectorConfig,
		siteRegionMapping,
		elastic,
		l,
	)

	hdl := handler.NewSolarmanCollectorHandler(serv)
	sub := api.Group("/collector")
	sub.GET("", hdl.Run)
}
