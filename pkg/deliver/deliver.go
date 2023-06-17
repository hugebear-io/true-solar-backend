package deliver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Payload struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func ResponseOK(ctx *gin.Context, data interface{}) {
	payload := Payload{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	}

	ctx.AbortWithStatusJSON(http.StatusOK, payload)
}

func ResponseCreated(ctx *gin.Context) {
	payload := Payload{
		Code:    http.StatusCreated,
		Message: "created",
	}

	ctx.AbortWithStatusJSON(http.StatusCreated, payload)
}

func ResponseUnauthorized(ctx *gin.Context) {
	payload := Payload{
		Code:    http.StatusUnauthorized,
		Message: "unauthorized",
	}

	ctx.AbortWithStatusJSON(http.StatusUnauthorized, payload)
}

func ResponseInternalError(ctx *gin.Context) {
	payload := Payload{
		Code:    http.StatusInternalServerError,
		Message: "system error",
	}

	ctx.AbortWithStatusJSON(http.StatusInternalServerError, payload)
}

func ResponseBadRequest(ctx *gin.Context, msg string) {
	payload := Payload{
		Code:    http.StatusBadRequest,
		Message: msg,
	}

	ctx.AbortWithStatusJSON(http.StatusBadRequest, payload)
}
