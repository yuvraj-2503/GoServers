package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"user-server/common"
	"user-server/signup/api"
	"user-server/signup/service"
	"user-server/validators"
)

type SignUpHandler struct {
	signupManager service.SignUpManager
}

func NewSignUpHandler(signupManager service.SignUpManager) *SignUpHandler {
	return &SignUpHandler{signupManager: signupManager}
}

func (s *SignUpHandler) SignUp(ctx *gin.Context) {
	var request api.SignUpRequest
	if !validators.ParseAndValidate(ctx, &request) {
		return
	}

	requestContext := ctx.Request.Context()
	result, err := s.signupManager.SignUp(&requestContext, &request)
	if err != nil {
		log.Printf("Failed to signup with sessionId %s , reason %v", request.SessionId, err)
		handle(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func handle(err error, ctx *gin.Context) {
	var userExistError *common.AlreadyExistsError
	if errors.As(err, &userExistError) {
		common.ConflictError(ctx, "user already exist")
		return
	}
	if common.HandleOtpErrors(ctx, err) {
		return
	}

	common.InternalError(ctx, "failed to process request, reason: "+err.Error())
	return
}
