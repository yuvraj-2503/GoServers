package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"otp-manager/common"
	userCommon "user-server/common"
	"user-server/validators"
)

func (a *AuthHandler) SendOTP(c *gin.Context) {
	var contact common.Contact
	if !validators.ParseAndValidate(c, &contact) {
		return
	}

	requestCtx := getRequestContext(c)
	sessionId, err := a.authService.SendOTP(&requestCtx, &contact)
	if err != nil {
		var notFoundErr *userCommon.NotFoundError
		switch {
		case errors.As(err, &notFoundErr):
			userCommon.NotFound(c, fmt.Sprintf("Could not send otp , reason : %s", err.Error()))
			return
		default:
			log.Printf("Error while sending otp , reason : %s", err.Error())
			userCommon.InternalError(c, fmt.Sprintf("Error while sending otp, reason : %s", err.Error()))
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"sessionId": sessionId})
}
