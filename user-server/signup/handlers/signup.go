package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"user-server/signup/api"
	"user-server/signup/service"
)

type SignUpHandler struct {
	signupManager service.SignUpManager
}

func NewSignUpHandler(signupManager service.SignUpManager) *SignUpHandler {
	return &SignUpHandler{signupManager: signupManager}
}

func (s *SignUpHandler) SignUp(ctx *gin.Context) {
	var request api.SignUpRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		return
	}

	requestContext := ctx.Request.Context()
	result, err := s.signupManager.SignUp(&requestContext, &request)
	if err != nil {
		log.Println(err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}
