package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
	"github.com/hugebear-io/true-solar-backend/pkg/deliver"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
)

type siteRegionMappingHandler struct {
	accessLogServ         domain.AccessLogService
	siteRegionMappingServ domain.SiteRegionMappingService
	logger                logger.Logger
}

func NewSiteRegionMappingHandler(
	accessLogServ domain.AccessLogService,
	siteRegionMappingServ domain.SiteRegionMappingService,
	logger logger.Logger,
) *siteRegionMappingHandler {
	return &siteRegionMappingHandler{
		accessLogServ:         accessLogServ,
		siteRegionMappingServ: siteRegionMappingServ,
		logger:                logger,
	}
}

func (api siteRegionMappingHandler) GetCity(ctx *gin.Context) {
	validate := validator.New()

	limitString := ctx.Query("limit")
	if err := validate.Var(limitString, "required,numeric,min=0"); err != nil {
		api.logger.Errorf("siteRegionMappingHandler.GetCity() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_PARAMETER)
		return
	}
	limit, _ := strconv.Atoi(limitString)

	offsetString := ctx.Query("offset")
	if err := validate.Var(offsetString, "required,numeric,min=0"); err != nil {
		api.logger.Errorf("siteRegionMappingHandler.GetCity() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_PARAMETER)
		return
	}
	offset, _ := strconv.Atoi(offsetString)

	if limit == -1 && offset == 0 {
		siteRegionMapping, err := api.siteRegionMappingServ.GetAllSiteRegionMapping()
		if err != nil {
			deliver.ResponseInternalError(ctx)
			return
		}

		result := map[string]interface{}{
			"total": len(siteRegionMapping),
			"data":  siteRegionMapping,
		}
		deliver.ResponseOK(ctx, result)
		return
	}

	siteRegionMapping, total, err := api.siteRegionMappingServ.GetSiteRegionMappingWithPagination(limit, offset)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	result := map[string]interface{}{
		"total": total,
		"data":  siteRegionMapping,
	}
	deliver.ResponseOK(ctx, result)
}

func (api siteRegionMappingHandler) GetRegion(ctx *gin.Context) {
	siteRegionMapping, err := api.siteRegionMappingServ.GetRegion()
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, siteRegionMapping)
}

func (api siteRegionMappingHandler) CreateCity(ctx *gin.Context) {
	var siteRegionMapping port.SiteRegionMapping
	if err := ctx.ShouldBindJSON(&siteRegionMapping); err != nil {
		api.logger.Errorf("siteRegionMappingHandler.CreateCity() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
		return
	}

	err := api.siteRegionMappingServ.CreateCity(siteRegionMapping)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	userID, exits := ctx.Get("user_id")
	if !exits {
		api.logger.Errorf("siteRegionMappingHandler.CreateCity() : %s", constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		return
	}

	err = api.accessLogServ.CreateAccessLog(port.AccessLog{
		ByUserID: userID.(int),
		Message:  "Create a new city",
	})

	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseCreated(ctx)
}

func (api siteRegionMappingHandler) UpdateCity(ctx *gin.Context) {
	id := ctx.Param("city_id")
	cityID, err := strconv.Atoi(id)
	if err != nil {
		api.logger.Errorf("siteRegionMappingHandler.UpdateCity() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_PARAMETER)
		return
	}

	var siteRegionMapping port.SiteRegionMapping
	if err := ctx.ShouldBindJSON(&siteRegionMapping); err != nil {
		api.logger.Errorf("siteRegionMappingHandler.UpdateCity() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
		return
	}

	err = api.siteRegionMappingServ.UpdateCity(cityID, siteRegionMapping)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	userID, exits := ctx.Get("user_id")
	if !exits {
		api.logger.Errorf("siteRegionMappingHandler.UpdateCity() : %s", constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		return
	}

	err = api.accessLogServ.CreateAccessLog(port.AccessLog{
		ByUserID: userID.(int),
		Message:  fmt.Sprintf("City (id : %d) updated", cityID),
	})
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, nil)
}

func (api siteRegionMappingHandler) UpdateRegion(ctx *gin.Context) {
	var regions domain.UpdateRegionsRequestBody
	if err := ctx.ShouldBindJSON(&regions); err != nil {
		api.logger.Errorf("siteRegionMappingHandler.UpdateRegion() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
		return
	}

	err := api.siteRegionMappingServ.UpdateRegion(regions)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	userID, exits := ctx.Get("user_id")
	if !exits {
		api.logger.Errorf("siteRegionMappingHandler.UpdateRegion() : %s", constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		return
	}

	err = api.accessLogServ.CreateAccessLog(port.AccessLog{
		ByUserID: userID.(int),
		Message:  fmt.Sprintf("Site region mapping (area : %s) updated", regions.Area),
	})
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, nil)
}

func (api siteRegionMappingHandler) DeleteCity(ctx *gin.Context) {
	id := ctx.Param("city_id")
	cityID, err := strconv.Atoi(id)
	if err != nil {
		api.logger.Errorf("siteRegionMappingHandler.DeleteCity() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_PARAMETER)
		return
	}

	err = api.siteRegionMappingServ.DeleteCity(cityID)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	userID, exits := ctx.Get("user_id")
	if !exits {
		api.logger.Errorf("siteRegionMappingHandler.DeleteCity() : %s", constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		return
	}

	err = api.accessLogServ.CreateAccessLog(port.AccessLog{
		ByUserID: userID.(int),
		Message:  fmt.Sprintf("City (id : %d) deleted", cityID),
	})
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, nil)
}

func (api siteRegionMappingHandler) DeleteArea(ctx *gin.Context) {
	area := ctx.Param("area")

	err := api.siteRegionMappingServ.DeleteArea(area)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	userID, exits := ctx.Get("user_id")
	if !exits {
		api.logger.Errorf("siteRegionMappingHandler.DeleteArea() : %s", constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		return
	}

	err = api.accessLogServ.CreateAccessLog(port.AccessLog{
		ByUserID: userID.(int),
		Message:  fmt.Sprintf("Area (area : %s) deleted", area),
	})
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, nil)
}

func (api siteRegionMappingHandler) ImportBMASiteRegionMapping(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		api.logger.Errorf("siteRegionMappingHandler.ImportBMASiteRegionMapping() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_PARAMETER)
		return
	}

	filename := file.Filename
	tempPath := filepath.Join(".", "temp")
	saveAs := filepath.Join(tempPath, filename)

	_ = os.MkdirAll(tempPath, os.ModePerm)

	err = ctx.SaveUploadedFile(file, saveAs)
	if err != nil {
		api.logger.Errorf("siteRegionMappingHandler.ImportBMASiteRegionMapping() : %s", err.Error())
		deliver.ResponseInternalError(ctx)
		return
	}

	defer os.Remove(saveAs)

	err = api.siteRegionMappingServ.ImportBMASiteRegionMapping(saveAs)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	userID, exits := ctx.Get("user_id")
	if !exits {
		api.logger.Errorf("siteRegionMappingHandler.ImportBMASiteRegionMapping() : %s", constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_CONTEXT_USER_ID_NOT_FOUND)
		return
	}

	err = api.accessLogServ.CreateAccessLog(port.AccessLog{
		ByUserID: userID.(int),
		Message:  fmt.Sprintf("Import BMA site mapping file (%s)", filename),
	})
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, nil)
}
