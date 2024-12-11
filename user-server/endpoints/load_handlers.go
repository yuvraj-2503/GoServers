package endpoints

import (
	"github.com/gin-gonic/gin"
	"user-server/common"
	"user-server/config"
	"user-server/endpoints/db"
	"user-server/endpoints/handlers"
	"user-server/endpoints/service"
)

var handler *handlers.UrlHandler

func LoadHandlers(router *gin.Engine) {
	mongoConfig := config.Configuration.MongoConfig
	urlColl, _ := mongoConfig.GetCollection(common.UrlsCollection)
	urlStore := db.NewUrlMongoStore(urlColl)
	urlService := service.NewUrlService(urlStore)
	handler = handlers.NewUrlHandler(urlService)
	loadRoutes(router)
}

func loadRoutes(router *gin.Engine) {
	group := router.Group("/api/v1")
	group.PATCH("/endpoints", handler.Upsert)
	group.GET("/endpoints", handler.GetAll)
}
