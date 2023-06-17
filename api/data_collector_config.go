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

func BindDataCollectorConfigAPI(router *gin.RouterGroup, db *sql.DB, logger logger.Logger) {
	alr := repo.NewAccessLogRepo(db)
	als := service.NewAccessLogService(alr)

	rep := repo.NewDataCollectorConfigRepo(db)
	serv := service.NewDataCollectorConfigService(rep)
	hdl := handler.NewDataCollectorConfigHandler(als, serv, logger)
	router.Group("/configs/collectors").
		Use(middleware.Authentication()).
		GET("", hdl.GetAllDataCollectorConfig).
		POST("", hdl.CreateDataCollectorConfig).
		GET("/:config_id", hdl.GetOneDataCollectorConfig).
		PUT("/:config_id", hdl.UpdateDataCollectorConfig).
		DELETE("/:config_id", hdl.DeleteDataCollectorConfig)
}
