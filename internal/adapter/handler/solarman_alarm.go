package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/pkg/deliver"
)

type solarmanAlarmHandler struct {
	serv domain.SolarmanAlarmService
}

func NewSolarmanAlarmHandler(serv domain.SolarmanAlarmService) *solarmanAlarmHandler {
	return &solarmanAlarmHandler{serv: serv}
}

func (h solarmanAlarmHandler) Run(c *gin.Context) {
	deliver.ResponseOK(c, nil)
	go h.serv.Run()
}
