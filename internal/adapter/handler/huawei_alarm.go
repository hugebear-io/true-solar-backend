package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/pkg/deliver"
)

type huaweiAlarmHandler struct {
	serv domain.HuaweiAlarmService
}

func NewHuaweiAlarmHandler(serv domain.HuaweiAlarmService) *huaweiAlarmHandler {
	return &huaweiAlarmHandler{serv: serv}
}

func (h huaweiAlarmHandler) Run(c *gin.Context) {
	deliver.ResponseOK(c, nil)
	go h.serv.Run()
}
