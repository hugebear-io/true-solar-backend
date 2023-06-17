package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/adapter/handler"
	"github.com/hugebear-io/true-solar-backend/internal/adapter/repo"
	"github.com/hugebear-io/true-solar-backend/internal/core/service"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
	"github.com/hugebear-io/true-solar-backend/pkg/middleware"
)

func BindSiteRegionMappingAPI(router *gin.RouterGroup, db *sql.DB, logger logger.Logger) {
	alr := repo.NewAccessLogRepo(db)
	als := service.NewAccessLogService(alr)

	rep := repo.NewSiteRegionMappingRepo(db)
	serv := service.NewSiteRegionMappingService(rep)
	hdl := handler.NewSiteRegionMappingHandler(als, serv, logger)
	router.Group("").
		Use(middleware.Authentication()).
		GET("/regions", hdl.GetRegion).
		PUT("/regions", hdl.UpdateRegion).
		POST("/regions/bma_upload", hdl.ImportBMASiteRegionMapping).
		GET("/cities", hdl.GetCity).
		POST("/cities", hdl.CreateCity).
		PUT("/cities/:city_id", hdl.UpdateCity).
		DELETE("/cities/:city_id", hdl.DeleteCity).
		DELETE("/areas/:area", hdl.DeleteArea)
}
