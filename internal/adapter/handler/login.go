package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/internal/core/port"
	"github.com/hugebear-io/true-solar-backend/pkg/constant"
	"github.com/hugebear-io/true-solar-backend/pkg/deliver"
	"github.com/hugebear-io/true-solar-backend/pkg/logger"
)

type loginHandler struct {
	accessLogServ domain.AccessLogService
	loginServ     domain.LoginService
	logger        logger.Logger
}

func NewLoginHandler(
	accessLogServ domain.AccessLogService,
	loginServ domain.LoginService,
	logger logger.Logger,
) *loginHandler {
	return &loginHandler{
		accessLogServ: accessLogServ,
		loginServ:     loginServ,
		logger:        logger,
	}
}

func (api loginHandler) Login(ctx *gin.Context) {
	var credential domain.Credential
	if err := ctx.ShouldBindJSON(&credential); err != nil {
		api.logger.Errorf("loginHandler.Login() : %s", err.Error())
		deliver.ResponseBadRequest(ctx, constant.RESPONSE_MESSAGE_INVALID_REQUEST_BODY)
		return
	}

	authToken, userID, err := api.loginServ.Login(credential)
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	err = api.accessLogServ.CreateAccessLog(port.AccessLog{
		ByUserID: userID,
		Message:  fmt.Sprintf("User (id : %d) logged in into system", userID),
	})

	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, authToken)
}
