package profile

import (
	blob_manager "blob-manager"
	"blob-manager/aws"
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
	var profileService = service.NewProfileService(profileStore, getBlobManager())
	profileHandler = handlers.NewProfileHandler(profileService)
	authHandler = auth.NewAuthHandler(config.Configuration.SecretKey)
	loadRoutes(router)
}

func getBlobManager() *blob_manager.BlobManager {
	var blobStore blob_manager.BlobStore
	blobStore = aws.CreateS3Client(getAwsConfig())
	blobManager := &blob_manager.BlobManager{
		BlobStore: blobStore,
	}
	return blobManager
}

func getAwsConfig() *aws.AWSConfig {
	return config.Configuration.AWSConfig
}

func loadRoutes(router *gin.Engine) {
	group := router.Group("/api/v1")
	group.Use(authHandler.Handle())
	{
		group.PATCH("/profile", profileHandler.UpdateProfile)
		group.GET("/profile", profileHandler.GetProfile)
		group.GET("/profile/:time", profileHandler.GetProfileByUserId)
	}
}
