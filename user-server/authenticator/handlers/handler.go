package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"user-server/authenticator/service"
)

type AuthHandler struct {
	authService service.UserAuthenticator
}

func NewAuthHandler(authService service.UserAuthenticator) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func getRequestContext(ctx *gin.Context) context.Context {
	return ctx.Request.Context()
}
