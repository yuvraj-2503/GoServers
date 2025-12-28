package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"user-server/authenticator"
	"user-server/config"
	"user-server/endpoints"
	endpointsdb "user-server/endpoints/db"
	"user-server/profile"
	profileDb "user-server/profile/db"
	"user-server/signin"
	"user-server/signup"
)

func main() {
	router := gin.Default()
	ctx := context.Background()

	// DB + handlers
	signup.LoadDB(&ctx)
	signup.LoadHandlers(router)

	signin.LoadHandlers(router)

	authenticator.LoadDB(&ctx)
	authenticator.LoadHandlers(router)

	profileDb.LoadDB(&ctx)
	profile.LoadHandlers(router)

	endpointsdb.LoadDB(&ctx)
	endpoints.LoadHandlers(router)

	// Health
	public := router.Group("/api/v1")
	public.GET("/health", Health)

	// PORT handling (Railway compatible)
	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.Itoa(config.Configuration.ServerPort)
	}

	log.Println("🚀 Starting user server on port:", port)

	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Panicf("Failed to start user server, reason: %v", err)
	}
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
