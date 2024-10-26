package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"user-server/common"
	"user-server/signup/api"
	"user-server/validators"
)

func (s *SignUpHandler) SendSmsOtp(ctx *gin.Context) {
	var request common.PhoneNumber
	if !validators.ParseAndValidate(ctx, &request) {
		return
	}

	if !validators.ValidatePhoneNumber(ctx, &request) {
		return
	}

	requestContext := ctx.Request.Context()
	result, err := s.signupManager.SendSmsOtp(&requestContext, &request)
	if err != nil {
		log.Printf("Failed to send sms otp to %s, reason %s",
			request.CountryCode+request.Number, err)
		handleError(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"sessionId": result})
}

func (s *SignUpHandler) SendEmailOtp(ctx *gin.Context) {
	var request api.EmailOTPRequest
	if !validators.ParseAndValidate(ctx, &request) {
		return
	}

	requestContext := ctx.Request.Context()
	result, err := s.signupManager.SendEmailOtp(&requestContext, request.EmailId)
	if err != nil {
		log.Printf("Failed to send email otp to %s, reason: %s", request.EmailId, err.Error())
		handleError(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"sessionId": result})
}

func handleError(err error, ctx *gin.Context) {
	var userExistsError *common.AlreadyExistsError
	if errors.As(err, &userExistsError) {
		common.ConflictError(ctx, "user already exists")
	} else {
		common.InternalError(ctx, "Failed to send otp, reason: "+err.Error())
	}
}
