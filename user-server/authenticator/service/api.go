package service

import (
	"context"
	"otp-manager/common"
	"otp-manager/otp"
	user "user-server/common"
)

type UserAuthenticator interface {
	SendOTP(ctx *context.Context, contact *common.Contact) (*string, error)
	Verify(ctx *context.Context, sessionId string, otp uint32) (*VerifyResponse, error)
}

type VerifyResponse struct {
	UserId string `json:"userId"`
	Token  string `json:"token"`
}

type UserAuthenticatorImpl struct {
	smsOtpManager   otp.OtpManager
	emailOtpManager otp.OtpManager
}

func NewUserAuthenticator() *UserAuthenticatorImpl {
	return &UserAuthenticatorImpl{}
}

func (manager *UserAuthenticatorImpl) SendOTP(ctx *context.Context, contact *common.Contact) (*string, error) {
	if contact.PhoneNumber != nil {
		return manager.sendSmsOtp(ctx, contact.PhoneNumber)
	} else {
		return manager.sendEmailOtp(ctx, contact.EmailId)
	}
}

func (manager *UserAuthenticatorImpl) sendEmailOtp(ctx *context.Context, emailId string) (*string, error) {
	return manager.emailOtpManager.Send(ctx, &common.Contact{
		EmailId: emailId,
	})
}

func (manager *UserAuthenticatorImpl) sendSmsOtp(ctx *context.Context, phoneNumber *user.PhoneNumber) (*string, error) {
	return manager.smsOtpManager.Send(ctx, &common.Contact{
		PhoneNumber: phoneNumber,
	})
}

func (manager *UserAuthenticatorImpl) Verify(ctx *context.Context, sessionId string, otp uint32) (*VerifyResponse, error) {

}
