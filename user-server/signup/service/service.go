package service

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	otpCommon "otp-manager/common"
	"otp-manager/otp"
	"time"
	token "token-manager"
	"user-server/common"
	"user-server/signup/api"
	"user-server/signup/db"
)

type SignUpManager interface {
	SendEmailOtp(ctx *context.Context, emailId string) (*string, error)
	SendSmsOtp(ctx *context.Context, phoneNumber *common.PhoneNumber) (*string, error)
	SignUp(ctx *context.Context, signUpRequest *api.SignUpRequest) (*common.AuthenticatedUser, error)
}

type MongoSignupManager struct {
	userStore       db.UserStore
	emailOtpManager otp.OtpManager
	smsOtpManager   otp.OtpManager
	tokenManager    token.TokenManager
}

func NewMongoSignupManager(userStore db.UserStore,
	emailOtpManager otp.OtpManager, smsOtpManager otp.OtpManager, tokenManager token.TokenManager) *MongoSignupManager {
	return &MongoSignupManager{
		userStore:       userStore,
		emailOtpManager: emailOtpManager,
		smsOtpManager:   smsOtpManager,
		tokenManager:    tokenManager,
	}
}

func (m *MongoSignupManager) SendEmailOtp(ctx *context.Context, emailId string) (*string, error) {
	result, _ := m.userStore.CheckExists(ctx, db.Filter{
		Key:   db.EmailId,
		Value: emailId,
	})

	if result {
		return nil, &common.AlreadyExistsError{Message: "Email already registered"}
	}
	return m.emailOtpManager.Send(ctx, &otpCommon.Contact{EmailId: emailId})
}

func (m *MongoSignupManager) SendSmsOtp(ctx *context.Context, phoneNumber *common.PhoneNumber) (*string, error) {
	result, _ := m.userStore.CheckIfMobileExists(ctx, phoneNumber)

	if result {
		return nil, &common.AlreadyExistsError{Message: "Phone Number already registered"}
	}
	return m.smsOtpManager.Send(ctx, &otpCommon.Contact{PhoneNumber: phoneNumber})
}

func (m *MongoSignupManager) SignUp(ctx *context.Context, signUpRequest *api.SignUpRequest) (*common.AuthenticatedUser, error) {
	_, err := m.verifyOtp(ctx, signUpRequest.SessionId, signUpRequest.OTP)
	if err != nil {
		return nil, err
	}

	user := m.createUser(signUpRequest.EmailId, signUpRequest.PhoneNumber)

	errCh := make(chan error)
	defer close(errCh)

	var result common.AuthenticatedUser
	m.insertUser(ctx, user, errCh)
	m.createApiKey(&result, user.UserId, user.EmailId, signUpRequest.Device.FingerPrint, signUpRequest.App, errCh)

	for i := 0; i < 2; i++ {
		err := <-errCh
		if err != nil {
			return nil, err
		}
	}

	result.UserId = user.UserId
	result.EmailId = user.EmailId
	result.PhoneNumber = user.PhoneNumber
	return &result, err
}

func (m *MongoSignupManager) verifyOtp(ctx *context.Context, sessionId string, otp uint64) (*otpCommon.Contact, error) {
	return m.emailOtpManager.Verify(ctx, sessionId, otp)
}

func (m *MongoSignupManager) createUser(emailId string, phoneNumber *common.PhoneNumber) *db.User {
	return &db.User{
		UserId:      primitive.NewObjectID().Hex(),
		EmailId:     emailId,
		PhoneNumber: phoneNumber,
	}
}

func (m *MongoSignupManager) insertUser(ctx *context.Context, user *db.User, errCh chan<- error) {
	go func() {
		err := m.userStore.Insert(ctx, user)
		errCh <- err
	}()
}

func (m *MongoSignupManager) createApiKey(
	result *common.AuthenticatedUser,
	userId string, emailId string, machineId string, appId string, errCh chan<- error) {
	go func() {
		claims := getClaims(userId, emailId, machineId, appId)
		apiKey, err := m.tokenManager.Generate(claims)
		if err == nil {
			result.ApiKey = apiKey
		}
		errCh <- err
	}()
}

func getClaims(userId string, emailId string, machineId string, appId string) *token.TokenClaims {
	iat := time.Now()
	exp := iat.Add(24 * 7 * time.Hour)
	return &token.TokenClaims{
		UserId:    userId,
		EmailId:   emailId,
		MachineId: machineId,
		App:       appId,
		IAT:       &iat,
		EXP:       &exp,
		Kind:      "USER",
	}
}
