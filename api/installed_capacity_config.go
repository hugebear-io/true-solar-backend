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

func BindInstalledCapacityConfigAPI(router *gin.RouterGroup, db *sql.DB, logger logger.Logger) {
	alr := repo.NewAccessLogRepo(db)
	als := service.NewAccessLogService(alr)

	rep := repo.NewInstalledCapacityConfigRepo(db)
	serv := service.NewInstalledCapacityConfigService(rep)
	hdl := handler.NewInstalledCapacityConfigHandler(als, serv, logger)
	router.Group("/configs/installed_capacity").
		Use(middleware.Authentication()).
		GET("", hdl.GetInstalledCapacityConfig).
		PUT("", hdl.UpdateInstalledCapacityConfig)
}
