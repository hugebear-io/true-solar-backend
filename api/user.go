package api

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/adapter/handler"
	"github.com/hugebear-io/true-solar-backend/internal/adapter/repo"
	"github.com/hugebear-io/true-solar-backend/internal/core/service"
	"github.com/hugebear-io/true-solar-backend/pkg/middleware"
)

func BindUserAPI(router *gin.RouterGroup, db *sql.DB) {
	rep := repo.NewUserRepo(db)
	serv := service.NewUserService(rep)
	hdl := handler.NewUserHandler(serv)
	router.Group("/kibana_auth").
		Use(middleware.Authentication()).
		GET("", hdl.GetKibanaUser)
}
