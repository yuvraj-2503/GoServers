package handlers

import (
	blobError "blob-manager/common"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
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
		var uploadErr *blobError.UploadError
		if errors.As(err, &uploadErr) {
			common.InternalError(ctx, "Profile Picture Upload failed")
			return
		}
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
		var downloadErr *blobError.DownloadError
		if errors.As(err, &downloadErr) {
			common.InternalError(ctx, "Profile Picture Download failed")
			return
		}
		common.InternalError(ctx, "Failed to get profile, reason: "+err.Error())
		return
	}

	if profile == nil {
		ctx.Status(http.StatusNoContent)
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

func (p *ProfileHandler) GetProfileByUserId(ctx *gin.Context) {
	userId := getUserIdFromContext(ctx)

	requestCtx := ctx.Request.Context()
	providedTime := epochToTime(ctx.Param("time"))
	profile, err := p.profileService.GetProfile(&requestCtx, userId, providedTime)
	if err != nil {
		var notFoundErr *common.NotFoundError
		if errors.As(err, &notFoundErr) {
			common.NotFound(ctx, err.Error())
			return
		}

		var downloadErr *blobError.DownloadError
		if errors.As(err, &downloadErr) {
			common.InternalError(ctx, "Profile Picture Download failed")
			return
		}
		common.InternalError(ctx, "Failed to get profile, reason: "+err.Error())
		return
	}

	if profile == nil {
		ctx.Status(http.StatusNoContent)
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

func epochToTime(epoch string) time.Time {
	epochTime, _ := strconv.ParseInt(epoch, 10, 64)
	return time.UnixMilli(epochTime)
}

func getUserIdFromContext(ctx *gin.Context) string {
	user, exists := ctx.Get("user")
	if !exists {
		return ""
	}

	userId := user.(token.TokenClaims).UserId
	return userId
}
