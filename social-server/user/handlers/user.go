package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"social-server/user/db"
	"social-server/user/service"
	token "token-manager"
	"user-server/common"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (s *UserHandler) RegisterUser(ctx *gin.Context) {
	var request db.UserDetails

	if err := ctx.ShouldBindJSON(&request); err != nil {
		common.BadRequest(ctx, "bad-request", "failed to parse user details")
		return
	}

	userId := getUserIdFromContext(ctx)
	requestCtx := ctx.Request.Context()
	request.UserId = userId

	result := s.service.RegisterUser(&requestCtx, &request)
	ctx.JSON(http.StatusOK, result)
}

func getUserIdFromContext(ctx *gin.Context) string {
	user, exists := ctx.Get("user")
	if !exists {
		return ""
	}

	userId := user.(token.TokenClaims).UserId
	return userId
}

func handleError(ctx *gin.Context, err error) {
	var alreadyExist *common.AlreadyExistsError
	if errors.As(err, &alreadyExist) {
		common.ConflictError(ctx, err.Error())
	} else {
		common.InternalError(ctx, "failed to register user, reason: "+err.Error())
	}
}
