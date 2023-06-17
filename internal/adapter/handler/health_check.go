package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/pkg/deliver"
)

type healthCheckHandler struct {
}

func NewHealthCheckHandler() *healthCheckHandler {
	return &healthCheckHandler{}
}

func (api healthCheckHandler) Metric(ctx *gin.Context) {
	deliver.ResponseOK(ctx, nil)
}
