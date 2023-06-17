package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
	"github.com/hugebear-io/true-solar-backend/pkg/deliver"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
)

type snmpConfigHandler struct {
	accessLogServ  domain.AccessLogService
	snmpConfigServ domain.SNMPConfigService
	logger         logger.Logger
}

func NewSNMPConfigHandler(
	accessLogServ domain.AccessLogService,
	snmpConfigServ domain.SNMPConfigService,
	logger logger.Logger,
) *snmpConfigHandler {
	return &snmpConfigHandler{
		accessLogServ:  accessLogServ,
		snmpConfigServ: snmpConfigServ,
		logger:         logger,
	}
}

func (api snmpConfigHandler) GetSNMPConfig(ctx *gin.Context) {
	snmpConfig, err := api.snmpConfigServ.GetSNMPConfig()
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, snmpConfig)
}

func (api snmpConfigHandler) UpdateSNMPConfig(ctx *gin.Context) {
	var snmpConfig port.SNMPConfig
	if err := ctx.ShouldBindJSON(&snmpConfig); err != nil {
		api.logger.Errorf("snmpConfigHandler.UpdateSNMPConfig() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
		return
	}

	err := api.snmpConfigServ.UpdateSNMPConfig(snmpConfig)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	userID, exits := ctx.Get("user_id")
	if !exits {
		api.logger.Errorf("snmpConfigHandler.UpdateSNMPConfig() : %s", constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		return
	}

	err = api.accessLogServ.CreateAccessLog(port.AccessLog{
		ByUserID: userID.(int),
		Message:  "SNMP configuration updated",
	})
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, nil)
}
