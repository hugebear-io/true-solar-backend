package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
	"github.com/hugebear-io/true-solar-backend/pkg/deliver"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
)

type elasticSearchConfigHandler struct {
	logger                  logger.Logger
	accessLogServ           domain.AccessLogService
	elasticsearchConfigServ domain.ElasticSearchConfigService
}

func NewElasticSearchConfigHandler(

	logger logger.Logger,
	accessLogServ domain.AccessLogService,
	elasticsearchConfigServ domain.ElasticSearchConfigService,
) *elasticSearchConfigHandler {
	return &elasticSearchConfigHandler{
		logger:                  logger,
		accessLogServ:           accessLogServ,
		elasticsearchConfigServ: elasticsearchConfigServ,
	}
}

func (api elasticSearchConfigHandler) GetElasticsearchConfig(ctx *gin.Context) {
	elasticsearchConfig, err := api.elasticsearchConfigServ.GetElasticsearchConfig()
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, elasticsearchConfig)
}

func (api elasticSearchConfigHandler) UpdateElasticsearchConfig(ctx *gin.Context) {
	var elasticsearchConfig port.ElasticSearchConfig
	if err := ctx.ShouldBindJSON(&elasticsearchConfig); err != nil {
		api.logger.Errorf("elasticSearchConfigHandler.UpdateElasticsearchConfig() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
		return
	}

	err := api.elasticsearchConfigServ.UpdateElasticsearchConfig(elasticsearchConfig)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	userID, exits := ctx.Get("user_id")
	if !exits {
		api.logger.Errorf("elasticSearchConfigHandler.UpdateElasticsearchConfig() : %s", constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		return
	}

	err = api.accessLogServ.CreateAccessLog(port.AccessLog{
		ByUserID: userID.(int),
		Message:  "Elasticsearch configuration updated",
	})
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, nil)
}
