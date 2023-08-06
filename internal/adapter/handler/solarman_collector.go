package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/pkg/deliver"
)

type solarmanCollectorHandler struct {
	serv domain.SolarmanCollectorService
}

func NewSolarmanCollectorHandler(serv domain.SolarmanCollectorService) *solarmanCollectorHandler {
	return &solarmanCollectorHandler{serv: serv}
}

func (h solarmanCollectorHandler) Run(c *gin.Context) {
	deliver.ResponseOK(c, nil)
	go h.serv.Run()
}
