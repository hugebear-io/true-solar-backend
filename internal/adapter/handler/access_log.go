package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
	"github.com/hugebear-io/true-solar-backend/pkg/deliver"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
)

type accessLogHandler struct {
	serv   domain.AccessLogService
	logger logger.Logger
}

func NewAccessLogHandler(serv domain.AccessLogService, logger logger.Logger) *accessLogHandler {
	return &accessLogHandler{
		serv:   serv,
		logger: logger,
	}
}

func (api accessLogHandler) GetAccessLog(ctx *gin.Context) {
	validate := validator.New()
	limitString := ctx.Query("limit")
	if err := validate.Var(limitString, "required,numeric,min=0"); err != nil {
		api.logger.Errorf("AccessLogAPI.GetAccessLog() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_PARAMETER)
		return
	}

	offsetString := ctx.Query("offset")
	if err := validate.Var(offsetString, "required,numeric,min=0"); err != nil {
		api.logger.Errorf("AccessLogAPI.GetAccessLog() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_PARAMETER)
		return
	}

	limit, _ := strconv.Atoi(limitString)
	offset, _ := strconv.Atoi(offsetString)
	accessLog, total, err := api.serv.GetAccessLog(limit, offset)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	result := map[string]interface{}{
		"total": total,
		"data":  accessLog,
	}
	deliver.ResponseOK(ctx, result)
}
