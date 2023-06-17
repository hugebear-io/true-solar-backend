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

type dataCollectorConfigHandler struct {
	logger                  logger.Logger
	dataCollectorConfigServ domain.DataCollectorConfigService
	accessLogServ           domain.AccessLogService
}

func NewDataCollectorConfigHandler(
	accessLogServ domain.AccessLogService,
	dataCollectorConfigServ domain.DataCollectorConfigService,
	logger logger.Logger,
) *dataCollectorConfigHandler {
	return &dataCollectorConfigHandler{
		logger:                  logger,
		dataCollectorConfigServ: dataCollectorConfigServ,
		accessLogServ:           accessLogServ}
}

func (api dataCollectorConfigHandler) GetAllDataCollectorConfig(ctx *gin.Context) {
	dataCollectorConfig, err := api.dataCollectorConfigServ.GetAllDataCollectorConfig()
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, dataCollectorConfig)
}

func (api dataCollectorConfigHandler) GetOneDataCollectorConfig(ctx *gin.Context) {
	id := ctx.Param("config_id")
	dataCollectorConfigID, err := strconv.Atoi(id)
	if err != nil {
		api.logger.Errorf("dataCollectorConfigHandler.GetOneDataCollectorConfig() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_PARAMETER)
		return
	}

	dataCollectorConfig, err := api.dataCollectorConfigServ.GetOneDataCollectorConfig(dataCollectorConfigID)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	userID, exits := ctx.Get("user_id")
	if !exits {
		api.logger.Errorf("dataCollectorConfigHandler.GetOneDataCollectorConfig() : %s", constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		return
	}

	err = api.accessLogServ.CreateAccessLog(port.AccessLog{
		ByUserID: userID.(int),
		Message:  fmt.Sprintf("Get collector vendor configuration (id : %d)", dataCollectorConfigID),
	})
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, dataCollectorConfig)
}

func (api dataCollectorConfigHandler) CreateDataCollectorConfig(ctx *gin.Context) {
	var dataCollectorConfig port.VendorAccount
	if err := ctx.ShouldBindJSON(&dataCollectorConfig); err != nil {
		api.logger.Errorf("dataCollectorConfigHandler.CreateDataCollectorConfig() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
		return
	}

	switch dataCollectorConfig.VendorType {
	case constant.VENDOR_TYPE_GROWATT:
		if dataCollectorConfig.Token == nil || *dataCollectorConfig.Token == "" {
			api.logger.Errorf("dataCollectorConfigHandler.CreateDataCollectorConfig() : %s", constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
			deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
			return
		}
	case constant.VENDOR_TYPE_KSTAR:
		// no-op
	case constant.VENDOR_TYPE_INVT:
		if dataCollectorConfig.AppID == nil || dataCollectorConfig.AppSecret == nil || *dataCollectorConfig.AppID == "" || *dataCollectorConfig.AppSecret == "" {
			api.logger.Errorf("dataCollectorConfigHandler.CreateDataCollectorConfig() : %s", constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
			deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
			return
		}
	case constant.VENDOR_TYPE_HUAWEI:
		// no-op
	default:
		api.logger.Errorf("dataCollectorConfigHandler.CreateDataCollectorConfig() : %s", constant.RESPONSE_MESSAGE_INVALID_VENDOR_TYPE)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_VENDOR_TYPE)
		return
	}

	err := api.dataCollectorConfigServ.CreateDataCollectorConfig(dataCollectorConfig)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	userID, exits := ctx.Get("user_id")
	if !exits {
		api.logger.Errorf("dataCollectorConfigHandler.CreateDataCollectorConfig() : %s", constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		return
	}

	err = api.accessLogServ.CreateAccessLog(port.AccessLog{
		ByUserID: userID.(int),
		Message:  "Create a new collector vendor configuration",
	})
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseCreated(ctx)
}

func (api dataCollectorConfigHandler) UpdateDataCollectorConfig(ctx *gin.Context) {
	id := ctx.Param("config_id")
	dataCollectorConfigID, err := strconv.Atoi(id)
	if err != nil {
		api.logger.Errorf("dataCollectorConfigHandler.UpdateDataCollectorConfig() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_PARAMETER)
		return
	}

	var dataCollectorConfig port.VendorAccount
	if err := ctx.ShouldBindJSON(&dataCollectorConfig); err != nil {
		api.logger.Errorf("dataCollectorConfigHandler.UpdateDataCollectorConfig() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
		return
	}

	switch dataCollectorConfig.VendorType {
	case constant.VENDOR_TYPE_GROWATT:
		if dataCollectorConfig.Token == nil || *dataCollectorConfig.Token == "" {
			api.logger.Errorf("dataCollectorConfigHandler.UpdateDataCollectorConfig() : %s", constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
			deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
			return
		}
	case constant.VENDOR_TYPE_KSTAR:
		// no-op
	case constant.VENDOR_TYPE_INVT:
		if dataCollectorConfig.AppID == nil || dataCollectorConfig.AppSecret == nil || *dataCollectorConfig.AppID == "" || *dataCollectorConfig.AppSecret == "" {
			api.logger.Errorf("dataCollectorConfigHandler.UpdateDataCollectorConfig() : %s", constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
			deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
			return
		}
	case constant.VENDOR_TYPE_HUAWEI:
		// no-op
	default:
		api.logger.Errorf("dataCollectorConfigHandler.UpdateDataCollectorConfig() : %s", constant.RESPONSE_MESSAGE_INVALID_VENDOR_TYPE)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_VENDOR_TYPE)
		return
	}

	err = api.dataCollectorConfigServ.UpdateDataCollectorConfig(dataCollectorConfigID, dataCollectorConfig)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	userID, exits := ctx.Get("user_id")
	if !exits {
		api.logger.Errorf("dataCollectorConfigHandler.UpdateDataCollectorConfig() : %s", constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		return
	}

	err = api.accessLogServ.CreateAccessLog(port.AccessLog{
		ByUserID: userID.(int),
		Message:  fmt.Sprintf("Collector vendor configuration (id : %d) updated", dataCollectorConfigID),
	})
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, nil)
}

func (api dataCollectorConfigHandler) DeleteDataCollectorConfig(ctx *gin.Context) {
	id := ctx.Param("config_id")
	dataCollectorConfigID, err := strconv.Atoi(id)
	if err != nil {
		api.logger.Errorf("dataCollectorConfigHandler.DeleteDataCollectorConfig() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_PARAMETER)
		return
	}

	err = api.dataCollectorConfigServ.DeleteDataCollectorConfig(dataCollectorConfigID)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	userID, exits := ctx.Get("user_id")
	if !exits {
		api.logger.Errorf("dataCollectorConfigHandler.DeleteDataCollectorConfig() : %s", constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		return
	}

	err = api.accessLogServ.CreateAccessLog(port.AccessLog{
		ByUserID: userID.(int),
		Message:  fmt.Sprintf("Collector vendor configuration (id : %d) deleted", dataCollectorConfigID),
	})
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, nil)
}
