package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
	"github.com/hugebear-io/true-solar-backend/pkg/deliver"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
)

type redisConfigHandler struct {
	accessLogServ   domain.AccessLogService
	redisConfigServ domain.RedisConfigService
	logger          logger.Logger
}

func NewRedisConfigHandler(
	accessLogServ domain.AccessLogService,
	redisConfigServ domain.RedisConfigService,
	logger logger.Logger,
) *redisConfigHandler {
	return &redisConfigHandler{
		accessLogServ:   accessLogServ,
		redisConfigServ: redisConfigServ,
		logger:          logger,
	}
}

func (api redisConfigHandler) GetRedisConfig(ctx *gin.Context) {
	redisConfig, err := api.redisConfigServ.GetRedisConfig()
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, redisConfig)
}

func (api redisConfigHandler) UpdateRedisConfig(ctx *gin.Context) {
	var redisConfig port.RedisConfig
	if err := ctx.ShouldBindJSON(&redisConfig); err != nil {
		api.logger.Errorf("redisConfigHandler.UpdateRedisConfig() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
		return
	}

	err := api.redisConfigServ.UpdateRedisConfig(redisConfig)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	userID, exits := ctx.Get("user_id")
	if !exits {
		api.logger.Errorf("redisConfigHandler.UpdateRedisConfig() : %s", constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		return
	}

	err = api.accessLogServ.CreateAccessLog(port.AccessLog{
		ByUserID: userID.(int),
		Message:  "Redis configuration updated",
	})
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, nil)
}
