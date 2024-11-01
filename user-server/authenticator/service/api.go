package service

import (
	"context"
	"otp-manager/common"
	"otp-manager/otp"
	token "token-manager"
	"user-server/authenticator/sessiondb"
	"user-server/signup/db"
)

type UserAuthenticator interface {
	SendOTP(ctx *context.Context, contact *common.Contact) (*string, error)
	Verify(ctx *context.Context, sessionId string, otp uint64) (*VerifyResponse, error)
}

type VerifyResponse struct {
	UserId string `json:"userId"`
	Token  string `json:"token"`
}

type UserAuthenticatorImpl struct {
	smsOtpManager   otp.OtpManager
	emailOtpManager otp.OtpManager
	userStore       db.UserStore
	tokenManager    token.TokenManager
	sessionMapping  sessiondb.SessionRepository
}

func NewUserAuthenticator(
	smsOtpManager otp.OtpManager,
	emailOtpManager otp.OtpManager,
	userStore db.UserStore,
	tokenManager token.TokenManager,
	sessionMapping sessiondb.SessionRepository,
) *UserAuthenticatorImpl {
	return &UserAuthenticatorImpl{
		smsOtpManager:   smsOtpManager,
		emailOtpManager: emailOtpManager,
		userStore:       userStore,
		tokenManager:    tokenManager,
		sessionMapping:  sessionMapping,
	}
}
