package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
	"github.com/hugebear-io/true-solar-backend/pkg/deliver"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
)

type alarmConfigHandler struct {
	accessLogServ   domain.AccessLogService
	alarmConfigServ domain.AlarmConfigService
	logger          logger.Logger
}

func NewAlarmConfigHandler(
	accessLogServ domain.AccessLogService,
	alarmConfigServ domain.AlarmConfigService,
	logger logger.Logger,
) *alarmConfigHandler {
	return &alarmConfigHandler{alarmConfigServ: alarmConfigServ, accessLogServ: accessLogServ, logger: logger}
}

func (api alarmConfigHandler) GetAllAlarmConfig(ctx *gin.Context) {
	alarmConfig, err := api.alarmConfigServ.GetAllAlarmConfig()
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, alarmConfig)
}

func (api alarmConfigHandler) GetOneAlarmConfig(ctx *gin.Context) {
	id := ctx.Param("config_id")
	alarmConfigID, err := strconv.Atoi(id)
	if err != nil {
		api.logger.Errorf("alarmConfigHandler.GetOneAlarmConfig() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_PARAMETER)
		return
	}

	alarmConfig, err := api.alarmConfigServ.GetOneAlarmConfig(alarmConfigID)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	userID, exits := ctx.Get("user_id")
	if !exits {
		api.logger.Errorf("alarmConfigHandler.GetOneAlarmConfig() : %s", constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		return
	}

	err = api.accessLogServ.CreateAccessLog(port.AccessLog{
		ByUserID: userID.(int),
		Message:  fmt.Sprintf("Get alarm vendor configuration (id : %d)", alarmConfigID),
	})

	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, alarmConfig)
}

func (api alarmConfigHandler) CreateAlarmConfig(ctx *gin.Context) {
	var alarmConfig port.VendorAccount
	if err := ctx.ShouldBindJSON(&alarmConfig); err != nil {
		api.logger.Errorf("alarmConfigHandler.CreateAlarmConfig() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
		return
	}

	switch alarmConfig.VendorType {
	case constant.VENDOR_TYPE_GROWATT:
		if alarmConfig.Token == nil || *alarmConfig.Token == "" {
			api.logger.Errorf("alarmConfigHandler.CreateAlarmConfig() : %s", constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
			deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
			return
		}
	case constant.VENDOR_TYPE_KSTAR:
		// no-op
	case constant.VENDOR_TYPE_INVT:
		if alarmConfig.AppID == nil || alarmConfig.AppSecret == nil || *alarmConfig.AppID == "" || *alarmConfig.AppSecret == "" {
			api.logger.Errorf("alarmConfigHandler.CreateAlarmConfig() : %s", constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
			deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
			return
		}
	case constant.VENDOR_TYPE_HUAWEI:
		// no-op
	default:
		api.logger.Errorf("alarmConfigHandler.CreateAlarmConfig() : %s", constant.RESPONSE_MESSAGE_INVALID_VENDOR_TYPE)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_VENDOR_TYPE)
		return
	}

	err := api.alarmConfigServ.CreateAlarmConfig(alarmConfig)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	userID, exits := ctx.Get("user_id")
	if !exits {
		api.logger.Errorf("alarmConfigHandler.CreateAlarmConfig() : %s", constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		return
	}

	err = api.accessLogServ.CreateAccessLog(port.AccessLog{
		ByUserID: userID.(int),
		Message:  "Create a new alarm vendor configuration",
	})

	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseCreated(ctx)
}

func (api alarmConfigHandler) UpdateAlarmConfig(ctx *gin.Context) {
	id := ctx.Param("config_id")
	alarmConfigID, err := strconv.Atoi(id)
	if err != nil {
		api.logger.Errorf("alarmConfigHandler.UpdateAlarmConfig() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_PARAMETER)
		return
	}

	var alarmConfig port.VendorAccount
	if err := ctx.ShouldBindJSON(&alarmConfig); err != nil {
		api.logger.Errorf("alarmConfigHandler.UpdateAlarmConfig() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
		return
	}

	switch alarmConfig.VendorType {
	case constant.VENDOR_TYPE_GROWATT:
		if alarmConfig.Token == nil || *alarmConfig.Token == "" {
			api.logger.Errorf("alarmConfigHandler.UpdateAlarmConfig() : %s", constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
			deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
			return
		}
	case constant.VENDOR_TYPE_KSTAR:
		// no-op
	case constant.VENDOR_TYPE_INVT:
		if alarmConfig.AppID == nil || alarmConfig.AppSecret == nil || *alarmConfig.AppID == "" || *alarmConfig.AppSecret == "" {
			api.logger.Errorf("alarmConfigHandler.UpdateAlarmConfig() : %s", constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
			deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
			return
		}
	case constant.VENDOR_TYPE_HUAWEI:
		// no-op
	default:
		api.logger.Errorf("alarmConfigHandler.UpdateAlarmConfig() : %s", constant.RESPONSE_MESSAGE_INVALID_VENDOR_TYPE)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_VENDOR_TYPE)
		return
	}

	err = api.alarmConfigServ.UpdateAlarmConfig(alarmConfigID, alarmConfig)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	userID, exits := ctx.Get("user_id")
	if !exits {
		api.logger.Errorf("alarmConfigHandler.UpdateAlarmConfig() : %s", constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		return
	}

	err = api.accessLogServ.CreateAccessLog(port.AccessLog{
		ByUserID: userID.(int),
		Message:  fmt.Sprintf("Alarm vendor configuration (id : %d) updated", alarmConfigID),
	})

	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, nil)
}

func (api alarmConfigHandler) DeleteAlarmConfig(ctx *gin.Context) {
	id := ctx.Param("config_id")
	alarmConfigID, err := strconv.Atoi(id)
	if err != nil {
		api.logger.Errorf("alarmConfigHandler.DeleteAlarmConfig() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_PARAMETER)
		return
	}

	err = api.alarmConfigServ.DeleteAlarmConfig(alarmConfigID)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	userID, exits := ctx.Get("user_id")
	if !exits {
		api.logger.Errorf("alarmConfigHandler.DeleteAlarmConfig() : %s", constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		return
	}

	err = api.accessLogServ.CreateAccessLog(port.AccessLog{
		ByUserID: userID.(int),
		Message:  fmt.Sprintf("Alarm vendor configuration (id : %d) deleted", alarmConfigID),
	})
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, nil)
}
