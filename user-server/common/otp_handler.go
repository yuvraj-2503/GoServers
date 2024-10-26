package common

import (
	"errors"
	"github.com/gin-gonic/gin"
	otpErrors "otp-manager/errors"
)

func HandleOtpErrors(ctx *gin.Context, err error) bool {
	var e *otpErrors.OtpError
	ok := errors.As(err, &e)
	if ok {
		if e.Code == otpErrors.INCORRECT {
			Unauthorized(ctx, "otp-incorrect", "otp is incorrect")
			return true
		} else if e.Code == otpErrors.EXPIRED {
			Unauthorized(ctx, "otp-expired", "otp is expired")
			return true
		} else if e.Code == otpErrors.LIMIT_EXCEEDED {
			TooManyRequest(ctx, "otp limit exceeded")
			return true
		} else if e.Code == otpErrors.NOT_FOUND {
			NotFound(ctx, "otp doesn't exist")
			return true

		}
	}
	return false
}
