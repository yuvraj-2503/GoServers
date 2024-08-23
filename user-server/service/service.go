package service

import (
	"context"
	otpCommon "otp-manager/common"
	"otp-manager/otp"
	"user-server/common"
	"user-server/db"
)

type SignUpManager interface {
	SendEmailOtp(ctx *context.Context, emailId string) (*string, error)
}

type MongoSignupManager struct {
	userStore       db.UserStore
	emailOtpManager otp.OtpManager
}

func NewMongoSignupManager(userStore db.UserStore,
	emailOtpManager otp.OtpManager) *MongoSignupManager {
	return &MongoSignupManager{
		userStore:       userStore,
		emailOtpManager: emailOtpManager,
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
