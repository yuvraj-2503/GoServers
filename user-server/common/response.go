package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Result struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func TooManyRequest(ctx *gin.Context, message string) {
	response := &Result{
		Code:    "too-many-requests",
		Message: message,
	}
	ctx.AbortWithStatusJSON(http.StatusTooManyRequests, response)
}

func NotFound(c *gin.Context, message string) {
	response := &Result{
		Code:    "not-found",
		Message: message,
	}
	c.AbortWithStatusJSON(http.StatusNotFound, response)
}

func BadRequest(c *gin.Context, code, message string) {
	response := &Result{
		Code:    code,
		Message: message,
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, response)
}

func ConflictError(c *gin.Context, message string) {
	response := &Result{
		Code:    "status-conflict",
		Message: message,
	}
	c.JSON(http.StatusConflict, response)
}

func InternalError(c *gin.Context, message string) {
	response := &Result{
		Code:    "core-error",
		Message: message,
	}
	c.JSON(http.StatusInternalServerError, response)
}

func Unauthorized(ctx *gin.Context, code, message string) {
	response := &Result{
		Code:    code,
		Message: message,
	}
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
}
