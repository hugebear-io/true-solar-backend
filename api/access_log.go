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

func BindAccessLogAPI(router *gin.RouterGroup, db *sql.DB, logger logger.Logger) {
	rep := repo.NewAccessLogRepo(db)
	serv := service.NewAccessLogService(rep)
	hdl := handler.NewAccessLogHandler(serv, logger)
	router.Group("").
		Use(middleware.Authentication()).
		GET("/access_log", hdl.GetAccessLog)
}
