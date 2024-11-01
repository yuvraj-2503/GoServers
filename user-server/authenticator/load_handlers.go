package authenticator

import (
	"context"
	"github.com/gin-gonic/gin"
	"mail-sender/mail"
	"otp-manager/otp"
	"otp-manager/senders"
	"otp-manager/sms"
	token "token-manager"
	"user-server/authenticator/handlers"
	"user-server/authenticator/service"
	"user-server/authenticator/sessiondb"
	"user-server/common"
	"user-server/config"
	"user-server/signup/db"
)

var authHandler *handlers.AuthHandler

func LoadHandlers(router *gin.Engine) {
	ctx := context.Background()
	authService := service.NewUserAuthenticator(getSmsOtpManager(&ctx), getEmailOtpManager(&ctx), getUserStore(), getTokenManager(),
		getSessionMapping())
	authHandler = handlers.NewAuthHandler(authService)
	loadRoutes(router)
}

func getSessionMapping() sessiondb.SessionRepository {
	mongoConfig := config.Configuration.MongoConfig
	sessionColl, _ := mongoConfig.GetCollection(common.SessionUserIdCollection)
	return sessiondb.NewSessionMongoStore(sessionColl)
}

func getTokenManager() token.TokenManager {
	return token.NewJwtTokenManager(config.Configuration.SecretKey)
}

func getUserStore() db.UserStore {
	mongoConfig := config.Configuration.MongoConfig
	userColl, _ := mongoConfig.GetCollection(common.UserCollection)
	return db.NewMongoUserStore(userColl)
}

func getEmailOtpManager(ctx *context.Context) otp.OtpManager {
	return otp.NewMongoOtpManager(getOtpStore(), getEmailOtpSender(ctx))
}

func getEmailOtpSender(ctx *context.Context) senders.OtpSender {
	sendGridSender := mail.NewSendgridMailSender(config.Configuration.SendgridConfig.SenderId,
		config.Configuration.SendgridConfig.SendgridApiKey)
	return senders.NewMailOtpSender(ctx, sendGridSender)
}

func getSmsOtpSender(ctx *context.Context) senders.OtpSender {
	smsSender := sms.NewMSG91Sender(
		config.Configuration.Msg91Config.BaseUrl,
		config.Configuration.Msg91Config.AuthKey,
		config.Configuration.Msg91Config.TemplateId)
	return senders.NewSmsOtpSender(ctx, smsSender)
}

func getOtpStore() otp.OTPStore {
	mongoConfig := config.Configuration.MongoConfig
	otpColl, _ := mongoConfig.GetCollection(common.OtpCollection)
	return otp.NewMongoOTPStore(otpColl)
}

func getSmsOtpManager(ctx *context.Context) otp.OtpManager {
	return otp.NewMongoOtpManager(getOtpStore(), getSmsOtpSender(ctx))
}

func loadRoutes(router *gin.Engine) {
	routes := router.Group("/api/v1/auth")
	routes.POST("/send/otp", authHandler.SendOTP)
	routes.POST("/verify/otp", authHandler.Verify)
}
