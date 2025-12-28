package signup

import (
	"context"
	"github.com/gin-gonic/gin"
	"mail-sender/mail"
	"otp-manager/otp"
	"otp-manager/senders"
	"otp-manager/sms"
	token_manager "token-manager"
	"user-server/common"
	"user-server/config"
	device_store "user-server/devices/db"
	device_manager "user-server/devices/service"
	"user-server/signup/db"
	"user-server/signup/handlers"
	"user-server/signup/service"
)

var signUpHandler *handlers.SignUpHandler

func LoadHandlers(router *gin.Engine) {
	ctx := context.Background()
	mongoConfig := config.Configuration.MongoConfig
	userCollection, _ := mongoConfig.GetCollection(common.UserCollection)
	otpCollection, _ := mongoConfig.GetCollection(common.OtpCollection)
	userStore := db.NewMongoUserStore(userCollection)
	otpStore := otp.NewMongoOTPStore(otpCollection)
	smsSender := sms.NewMSG91Sender(
		config.Configuration.Msg91Config.BaseUrl,
		config.Configuration.Msg91Config.AuthKey,
		config.Configuration.Msg91Config.TemplateId)

	smsOtpSender := senders.NewSmsOtpSender(&ctx, smsSender)

	sendGridSender := mail.NewSendgridMailSender(config.Configuration.SendgridConfig.SenderId,
		config.Configuration.SendgridConfig.SendgridApiKey)
	mailOtpSender := senders.NewMailOtpSender(&ctx, sendGridSender)
	emailOtpManager := otp.NewMongoOtpManager(otpStore, mailOtpSender)
	smsOtpManager := otp.NewMongoOtpManager(otpStore, smsOtpSender)
	tokenManager := token_manager.NewJwtTokenManager(config.Configuration.SecretKey)
	deviceCollection, _ := mongoConfig.GetCollection(common.UserDeviceCollection)
	userDeviceStore := device_store.NewMongoUserDeviceStore(deviceCollection)
	deviceManager := device_manager.NewUserDeviceManager(userDeviceStore)

	signUpManager := service.NewMongoSignupManager(userStore, emailOtpManager, smsOtpManager, tokenManager, deviceManager)
	signUpHandler = handlers.NewSignUpHandler(signUpManager)
	loadRoutes(router)
}

func loadRoutes(router *gin.Engine) {
	signup := router.Group("/api/v1")
	signup.POST("/signup", signUpHandler.SignUp)
	signup.POST("/otp/sms", signUpHandler.SendSmsOtp)
	signup.POST("/otp/email", signUpHandler.SendEmailOtp)
}
