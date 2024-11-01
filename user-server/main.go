package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"user-server/authenticator"
	"user-server/config"
	"user-server/signin"
	"user-server/signup"
)

func main() {
	router := gin.Default()
	ctx := context.Background()

	// signup
	signup.LoadDB(&ctx)
	signup.LoadHandlers(router)

	// signin
	signin.LoadHandlers(router)

	// authenticator
	authenticator.LoadDB(&ctx)
	authenticator.LoadHandlers(router)

	public := router.Group("/api/v1")
	public.GET("/health", Health)

	err := router.Run("0.0.0.0:" + strconv.Itoa(config.Configuration.ServerPort))
	if err != nil {
		log.Panicf("Failed to start user server, reason: %v", err)
		return
	}

	log.Println("Started user server at port: " + strconv.Itoa(config.Configuration.ServerPort))
}

func Health(c *gin.Context) {
	c.Data(http.StatusOK, gin.MIMEPlain, []byte{'0'})
}
