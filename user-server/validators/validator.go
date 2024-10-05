package validators

import (
	"github.com/gin-gonic/gin"
	"log"
	"user-server/common"
)

func ParseAndValidate(ctx *gin.Context, object interface{}) bool {
	preCheck(ctx)
	if err := ctx.ShouldBindJSON(object); err != nil {
		log.Printf("Validation Error: %s", err) // TODO: We can handle this error to provide more fine grained message to the client
		common.BadRequest(ctx, "bad-request", "Invalid request: failed to parse request body")
		return false
	}
	return true
}

func preCheck(ctx *gin.Context) {
	if ctx.Request.ContentLength == 0 {
		common.BadRequest(ctx, "bad-request", "request body is empty")
		return
	}
}

func ValidatePhoneNumber(ctx *gin.Context, phone *common.PhoneNumber) {
	switch {
	case len(phone.Number) < 7 || len(phone.Number) > 15:
		common.BadRequest(ctx, "bad-request", "invalid mobile number")
	}
}
