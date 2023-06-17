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

func BindSNMPConfigAPI(router *gin.RouterGroup, db *sql.DB, logger logger.Logger) {
	alr := repo.NewAccessLogRepo(db)
	als := service.NewAccessLogService(alr)

	rep := repo.NewSNMPConfigRepo(db)
	serv := service.NewSNMPConfigService(rep)
	hdl := handler.NewSNMPConfigHandler(als, serv, logger)
	router.Group("/configs/snmp").
		Use(middleware.Authentication()).
		GET("", hdl.GetSNMPConfig).
		PUT("", hdl.UpdateSNMPConfig)
}
