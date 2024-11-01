package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"user-server/common"
	"user-server/validators"
)

func (a *AuthHandler) Verify(c *gin.Context) {
	var verifyRequest VerifyRequest
	if !validators.ParseAndValidate(c, &verifyRequest) {
		return
	}
	ctx := getRequestContext(c)
	verifyResponse, err := a.authService.Verify(&ctx, verifyRequest.SessionId, verifyRequest.Otp)
	if err != nil {
		log.Printf("Failed to verify otp with sessionId: %s, reason: %s", verifyRequest.SessionId, err.Error())
		handle(err, c)
		return
	}

	c.JSON(http.StatusOK, verifyResponse)
}

func handle(err error, ctx *gin.Context) {
	var sessionNotFoundErr *common.NotFoundError
	if errors.As(err, &sessionNotFoundErr) {
		common.NotFound(ctx, "session does not exist")
		return
	}
	if common.HandleOtpErrors(ctx, err) {
		return
	}

	common.InternalError(ctx, "failed to verify otp, reason: "+err.Error())
	return
}

type VerifyRequest struct {
	SessionId string `json:"sessionId"`
	Otp       uint64 `json:"otp"`
}
