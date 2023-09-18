package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/pkg/deliver"
)

type SolarmanAlarmHandler struct {
	serv domain.SolarmanAlarmService
}

func NewSolarmanAlarmHandler(serv domain.SolarmanAlarmService) *SolarmanAlarmHandler {
	return &SolarmanAlarmHandler{serv: serv}
}

func (h SolarmanAlarmHandler) Run(c *gin.Context) {
	deliver.ResponseOK(c, nil)
	go h.serv.Run()
}

func (h SolarmanAlarmHandler) RunJob() {
	go h.serv.Run()
}
