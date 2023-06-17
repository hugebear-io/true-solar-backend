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

func BindAlarmConfigAPI(router *gin.RouterGroup, db *sql.DB, logger logger.Logger) {
	alr := repo.NewAccessLogRepo(db)
	als := service.NewAccessLogService(alr)

	rep := repo.NewAlarmConfigRepo(db)
	serv := service.NewAlarmConfigService(rep)
	hdl := handler.NewAlarmConfigHandler(als, serv, logger)
	router.Group("/configs/alarms").
		Use(middleware.Authentication()).
		GET("", hdl.GetAllAlarmConfig).
		POST("", hdl.CreateAlarmConfig).
		GET("/:config_id", hdl.GetOneAlarmConfig).
		PUT("/:config_id", hdl.UpdateAlarmConfig).
		DELETE("/:config_id", hdl.DeleteAlarmConfig)
}
