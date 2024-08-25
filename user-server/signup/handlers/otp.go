package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"user-server/common"
)

func (s *SignUpHandler) SendSmsOtp(ctx *gin.Context) {
	var request common.PhoneNumber
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Println(err)
		return
	}

	requestContext := ctx.Request.Context()
	result, err := s.signupManager.SendSmsOtp(&requestContext, &request)
	if err != nil {
		log.Println(err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"sessionId": result})
}
