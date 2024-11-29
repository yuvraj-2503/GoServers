package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	token "token-manager"
	"user-server/common"
	"user-server/profile/service"
)

type ProfileHandler struct {
	profileService service.ProfileService
}

func NewProfileHandler(profileService service.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

func (p *ProfileHandler) UpdateProfile(ctx *gin.Context) {
	userId := getUserIdFromContext(ctx)
	var userProfile service.UserProfile
	if err := ctx.ShouldBindJSON(&userProfile); err != nil {
		common.BadRequest(ctx, "bad-request", "Failed to parse request body")
		return
	}

	requestCtx := ctx.Request.Context()
	err := p.profileService.UpsertProfile(&requestCtx, userId, &userProfile)
	if err != nil {
		common.InternalError(ctx, "Failed to update profile, reason: "+err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}

func (p *ProfileHandler) GetProfile(ctx *gin.Context) {
	userId := getUserIdFromContext(ctx)

	requestCtx := ctx.Request.Context()
	profile, err := p.profileService.GetProfileByUserId(&requestCtx, userId)
	if err != nil {
		var notFoundErr *common.NotFoundError
		if errors.As(err, &notFoundErr) {
			common.NotFound(ctx, err.Error())
			return
		}
		common.InternalError(ctx, "Failed to get profile, reason: "+err.Error())
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

func getUserIdFromContext(ctx *gin.Context) string {
	user, exists := ctx.Get("user")
	if !exists {
		return ""
	}

	userId := user.(token.TokenClaims).UserId
	return userId
}
