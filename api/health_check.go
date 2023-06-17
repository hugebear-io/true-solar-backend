package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/adapter/handler"
)

func BindHealthCheckAPI(api *gin.RouterGroup) {
	h := handler.NewHealthCheckHandler()
	api.GET("/metric", h.Metric)
}
