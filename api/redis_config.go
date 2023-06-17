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

func BindRedisConfigAPI(router *gin.RouterGroup, db *sql.DB, logger logger.Logger) {
	alr := repo.NewAccessLogRepo(db)
	als := service.NewAccessLogService(alr)

	rep := repo.NewRedisConfigRepo(db)
	serv := service.NewRedisConfigService(rep)
	hdl := handler.NewRedisConfigHandler(als, serv, logger)
	router.Group("/configs/redis").
		Use(middleware.Authentication()).
		GET("", hdl.GetRedisConfig).
		PUT("")
}
