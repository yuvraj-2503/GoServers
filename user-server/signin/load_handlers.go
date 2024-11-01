package signin

import (
	"github.com/gin-gonic/gin"
	token_manager "token-manager"
	"user-server/common"
	"user-server/config"
	"user-server/signin/handlers"
	"user-server/signin/service"
	"user-server/signup/db"
)

var signInHandler *handlers.SignInHandler

func LoadHandlers(router *gin.Engine) {
	mongoConfig := config.Configuration.MongoConfig
	userCollection, _ := mongoConfig.GetCollection(common.UserCollection)
	userStore := db.NewMongoUserStore(userCollection)
	tokenManager := token_manager.NewJwtTokenManager(config.Configuration.SecretKey)
	signInManager := service.NewMongoSignInManager(userStore, tokenManager)
	signInHandler = handlers.NewSignInHandler(signInManager)
	loadRoutes(router)
}

func loadRoutes(router *gin.Engine) {
	signin := router.Group("/api/v1")
	signin.POST("/signIn", signInHandler.SignIn)
}
