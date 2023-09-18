package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/pkg/deliver"
)

type SolarmanCollectorHandler struct {
	serv domain.SolarmanCollectorService
}

func NewSolarmanCollectorHandler(serv domain.SolarmanCollectorService) *SolarmanCollectorHandler {
	return &SolarmanCollectorHandler{serv: serv}
}

func (h SolarmanCollectorHandler) Run(c *gin.Context) {
	deliver.ResponseOK(c, nil)
	go h.serv.Run()
}

func (h SolarmanCollectorHandler) RunJob() {
	go h.serv.Run()
}
