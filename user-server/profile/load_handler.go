package profile

import (
	"github.com/gin-gonic/gin"
	"user-server/auth"
	"user-server/common"
	"user-server/config"
	"user-server/profile/db"
	"user-server/profile/handlers"
	"user-server/profile/service"
)

var profileHandler *handlers.ProfileHandler
var authHandler *auth.AuthHandler

func LoadHandlers(router *gin.Engine) {
	mongoConfig := config.Configuration.MongoConfig
	profileColl, _ := mongoConfig.GetCollection(common.ProfileCollection)
	var profileStore = db.NewMongoProfileStore(profileColl)
	var profileService = service.NewProfileService(profileStore)
	profileHandler = handlers.NewProfileHandler(profileService)
	authHandler = auth.NewAuthHandler(config.Configuration.SecretKey)
	loadRoutes(router)
}

func loadRoutes(router *gin.Engine) {
	group := router.Group("/api/v1")
	group.Use(authHandler.Handle())
	{
		group.POST("/profile", profileHandler.UpdateProfile)
		group.GET("/profile", profileHandler.GetProfile)
	}
}
