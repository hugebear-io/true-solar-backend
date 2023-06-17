package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
	"github.com/hugebear-io/true-solar-backend/pkg/deliver"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
)

type performanceAlarmConfigHandler struct {
	accessLogServ         domain.AccessLogService
	performanceConfigServ domain.PerformanceAlarmConfigService
	logger                logger.Logger
}

func NewPerformanceAlarmConfigHandler(
	accessLogServ domain.AccessLogService,
	performanceConfigServ domain.PerformanceAlarmConfigService,
	logger logger.Logger,
) *performanceAlarmConfigHandler {
	return &performanceAlarmConfigHandler{
		accessLogServ:         accessLogServ,
		performanceConfigServ: performanceConfigServ,
		logger:                logger,
	}
}

func (api performanceAlarmConfigHandler) GetPerformanceAlarmConfig(ctx *gin.Context) {
	performanceAlarmConfig, err := api.performanceConfigServ.GetPerformanceAlarmConfig()
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, performanceAlarmConfig)
}

func (api performanceAlarmConfigHandler) UpdatePerformanceAlarmConfig(ctx *gin.Context) {
	var performanceAlarmConfig domain.UpdatePerformanceAlarmRequestBody
	if err := ctx.ShouldBindJSON(&performanceAlarmConfig); err != nil {
		api.logger.Errorf("performanceAlarmConfigHandler.UpdatePerformanceAlarmConfig() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
		return
	}

	err := api.performanceConfigServ.UpdatePerformanceAlarmConfig(performanceAlarmConfig)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	userID, exits := ctx.Get("user_id")
	if !exits {
		api.logger.Errorf("performanceAlarmConfigHandler.UpdatePerformanceAlarmConfig() : %s", constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		return
	}

	err = api.accessLogServ.CreateAccessLog(port.AccessLog{
		ByUserID: userID.(int),
		Message:  "Performance alarm configuration updated",
	})

	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, nil)
}
