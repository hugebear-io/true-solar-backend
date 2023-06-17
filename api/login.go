package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/adapter/handler"
	"github.com/hugebear-io/true-solar-backend/internal/adapter/repo"
	"github.com/hugebear-io/true-solar-backend/internal/core/service"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
)

func BindLoginAPI(router *gin.RouterGroup, secretKey string, db *sql.DB, logger logger.Logger) {
	alr := repo.NewAccessLogRepo(db)
	als := service.NewAccessLogService(alr)

	rep := repo.NewUserRepo(db)
	serv := service.NewLoginService(rep, secretKey)
	hdl := handler.NewLoginHandler(als, serv, logger)
	router.Group("/login").
		POST("", hdl.Login)
}
