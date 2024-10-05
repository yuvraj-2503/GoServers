package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"user-server/signin/api"
	"user-server/signin/service"
)

type SignInHandler struct {
	signInManager service.SignInManager
}

func NewSignInHandler(signInManager service.SignInManager) *SignInHandler {
	return &SignInHandler{signInManager: signInManager}
}

func (s *SignInHandler) SignIn(ctx *gin.Context) {
	var request api.SignInRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		return
	}

	requestContext := ctx.Request.Context()
	result, err := s.signInManager.SignIn(&requestContext, &request)
	if err != nil {
		log.Println(err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}
