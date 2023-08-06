package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/adapter/handler"
	"github.com/hugebear-io/true-solar-backend/internal/adapter/repo"
	"github.com/hugebear-io/true-solar-backend/internal/core/service"
	"github.com/hugebear-io/true-solar-backend/internal/infra"
	"github.com/hugebear-io/true-solar-backend/pkg/config"
)

func BindSolarmanCollectorAPI(api *gin.RouterGroup) {
	elasticConfig := config.Config.ElasticSearch
	elastic := repo.NewElasticSearchRepo(infra.ElasticSearch, elasticConfig.Index)

	siteRegionMappingRepo := repo.NewSiteRegionMappingRepo(infra.SqlDB)
	siteRegionMapping := service.NewSiteRegionMappingService(siteRegionMappingRepo)

	dataCollectorConfigRepo := repo.NewDataCollectorConfigRepo(infra.SqlDB)
	dataCollectorConfig := service.NewDataCollectorConfigService(dataCollectorConfigRepo)

	serv := service.NewSolarmanCollectorService(
		dataCollectorConfig,
		siteRegionMapping,
		elastic,
	)

	hdl := handler.NewSolarmanCollectorHandler(serv)
	sub := api.Group("/collector")
	sub.GET("", hdl.Run)
}
