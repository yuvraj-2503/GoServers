package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"user-server/common"
	"user-server/endpoints/service"
)

type UrlHandler struct {
	service service.UrlService
}

func NewUrlHandler(service service.UrlService) *UrlHandler {
	return &UrlHandler{
		service: service,
	}
}

func (u *UrlHandler) Upsert(c *gin.Context) {
	ctx := c.Request.Context()
	env := c.Query("env")
	var urls map[string]string
	if err := c.ShouldBindJSON(&urls); err != nil {
		common.BadRequest(c, "Invalid-request-body", "failed to parse request body")
		return
	}

	err := u.service.Upsert(&ctx, urls, env)
	if err != nil {
		common.InternalError(c, "failed to upsert urls data, reason: "+err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (u *UrlHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	env := c.Query("env")
	result, err := u.service.GetAll(&ctx, env)
	if err != nil {
		var notFoundError *common.NotFoundError
		if errors.As(err, &notFoundError) {
			common.NotFound(c, "Urls not found")
			return
		}
		common.InternalError(c, "failed to get urls data, reason: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
