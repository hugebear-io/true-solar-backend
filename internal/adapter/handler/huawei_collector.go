package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/pkg/deliver"
)

type huaweiCollectorHandler struct {
	serv domain.HuaweiCollectorService
}

func NewHuaweiCollectorHandler(serv domain.HuaweiCollectorService) *huaweiCollectorHandler {
	return &huaweiCollectorHandler{serv: serv}
}

func (h huaweiCollectorHandler) Run(c *gin.Context) {
	deliver.ResponseOK(c, nil)
	go h.serv.Run()
}
