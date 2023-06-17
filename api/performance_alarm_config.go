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

func BindPerformanceAlarmConfigAPI(router *gin.RouterGroup, db *sql.DB, logger logger.Logger) {
	alr := repo.NewAccessLogRepo(db)
	als := service.NewAccessLogService(alr)

	rep := repo.NewPerformanceAlarmConfigRepo(db)
	serv := service.NewPerformanceAlarmConfigService(rep)
	hdl := handler.NewPerformanceAlarmConfigHandler(als, serv, logger)
	router.Group("/configs/performance").
		Use(middleware.Authentication()).
		GET("", hdl.GetPerformanceAlarmConfig).
		PUT("", hdl.UpdatePerformanceAlarmConfig)
}
