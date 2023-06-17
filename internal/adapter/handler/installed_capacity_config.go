package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
	"github.com/hugebear-io/true-solar-backend/pkg/deliver"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
)

type installedCapacityConfigHandler struct {
	accessLogServ               domain.AccessLogService
	installedCapacityConfigServ domain.InstalledCapacityConfigService
	logger                      logger.Logger
}

func NewInstalledCapacityConfigHandler(
	accessLogServ domain.AccessLogService,
	installedCapacityConfigServ domain.InstalledCapacityConfigService,
	logger logger.Logger,
) *installedCapacityConfigHandler {
	return &installedCapacityConfigHandler{
		accessLogServ:               accessLogServ,
		installedCapacityConfigServ: installedCapacityConfigServ,
		logger:                      logger,
	}
}

func (api installedCapacityConfigHandler) GetInstalledCapacityConfig(ctx *gin.Context) {
	installedCapacityConfig, err := api.installedCapacityConfigServ.GetInstalledCapacityConfig()
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, installedCapacityConfig)
}

func (api installedCapacityConfigHandler) UpdateInstalledCapacityConfig(ctx *gin.Context) {
	var installedCapacityConfig port.InstalledCapacityConfig
	if err := ctx.ShouldBindJSON(&installedCapacityConfig); err != nil {
		api.logger.Errorf("installedCapacityConfigHandler.UpdateInstallCapacityConfig() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
		return
	}

	err := api.installedCapacityConfigServ.UpdateInstalledCapacityConfig(installedCapacityConfig)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	userID, exits := ctx.Get("user_id")
	if !exits {
		api.logger.Errorf("installedCapacityConfigHandler.UpdateInstallCapacityConfig() : %s", constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		return
	}

	err = api.accessLogServ.CreateAccessLog(port.AccessLog{
		ByUserID: userID.(int),
		Message:  "Estimated system KWh configuration updated",
	})

	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, nil)
}
