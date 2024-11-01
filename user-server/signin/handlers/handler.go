package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"user-server/common"
	"user-server/signin/api"
	"user-server/signin/service"
	"user-server/validators"
)

type SignInHandler struct {
	signInManager service.SignInManager
}

func NewSignInHandler(signInManager service.SignInManager) *SignInHandler {
	return &SignInHandler{signInManager: signInManager}
}

func (s *SignInHandler) SignIn(ctx *gin.Context) {
	var request api.SignInRequest
	if !validators.ParseAndValidate(ctx, &request) {
		return
	}

	requestContext := ctx.Request.Context()
	result, err := s.signInManager.SignIn(&requestContext, &request)
	if err != nil {
		handleSignInErrors(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func handleSignInErrors(ctx *gin.Context, err error) {
	var userNotExistError *common.UserDoesNotExistError
	if errors.As(err, &userNotExistError) {
		common.NotFound(ctx, "user doesn't exist")
	} else {
		log.Printf("failed to signin, reason: %s", err.Error())
		common.InternalError(ctx, fmt.Sprintf("failed to sign in, reason: %s", err.Error()))
	}
}
