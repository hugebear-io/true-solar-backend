package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hugebear-io/true-solar-backend/internal/core/domain"
	"github.com/hugebear-io/true-solar-backend/pkg/deliver"
)

type userHandler struct {
	userServ domain.UserService
}

func NewUserHandler(userServ domain.UserService) *userHandler {
	return &userHandler{userServ: userServ}
}

func (api userHandler) GetKibanaUser(ctx *gin.Context) {
	user, err := api.userServ.GetKibanaUser()
	if err != nil {
		deliver.ResponseInternalError(ctx)
		return
	}

	deliver.ResponseOK(ctx, user)
}
