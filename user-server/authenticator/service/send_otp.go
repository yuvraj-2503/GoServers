package service

import (
	"context"
	"otp-manager/common"
	"time"
	"user-server/authenticator/sessiondb"
	userCommon "user-server/common"
	"user-server/signup/db"
)

func (manager *UserAuthenticatorImpl) SendOTP(ctx *context.Context, contact *common.Contact) (*string, error) {
	if contact.PhoneNumber != nil {
		return manager.sendSmsOtp(ctx, contact)
	} else {
		return manager.sendEmailOtp(ctx, contact)
	}
}

func (manager *UserAuthenticatorImpl) sendEmailOtp(ctx *context.Context, contact *common.Contact) (*string, error) {
	userId, err := manager.getUserIdFromEmail(ctx, contact.EmailId)
	if err == nil {
		if userId != nil {
			sessionId, err := manager.emailOtpManager.Send(ctx, contact)
			if err != nil {
				return nil, err
			}

			err = manager.storeUserIdSessionId(ctx, userId, sessionId)
			if err == nil {
				return sessionId, err
			}
			return nil, err
		}
		return nil, &userCommon.NotFoundError{Message: "user not found for emailId: " + contact.EmailId}
	}

	return nil, err
}

func (manager *UserAuthenticatorImpl) sendSmsOtp(ctx *context.Context, contact *common.Contact) (*string, error) {
	userId, err := manager.getUserIdFromPhone(ctx, contact.PhoneNumber)
	if err == nil {
		if userId != nil {
			sessionId, err := manager.smsOtpManager.Send(ctx, contact)
			if err != nil {
				return nil, err
			}

			err = manager.storeUserIdSessionId(ctx, userId, sessionId)
			if err == nil {
				return sessionId, err
			}
			return nil, err
		}

		return nil, &userCommon.NotFoundError{Message: "user not found for phoneNumber: " +
			contact.PhoneNumber.CountryCode + contact.PhoneNumber.Number}
	}
	return nil, err
}

func (manager *UserAuthenticatorImpl) getUserIdFromEmail(ctx *context.Context, emailId string) (*string, error) {
	user, err := manager.userStore.Get(ctx, db.Filter{
		Key:   db.EmailId,
		Value: emailId,
	})
	if err != nil {
		return nil, err
	}
	return &user.UserId, nil
}

func (manager *UserAuthenticatorImpl) getUserIdFromPhone(ctx *context.Context, number *userCommon.PhoneNumber) (*string, error) {
	user, err := manager.userStore.GetByPhoneNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	return &user.UserId, nil
}

func (manager *UserAuthenticatorImpl) storeUserIdSessionId(ctx *context.Context, userId, sessionId *string) error {
	var now = time.Now()
	err := manager.sessionMapping.Insert(ctx, &sessiondb.SessionMapping{
		SessionId: *sessionId,
		UserId:    *userId,
		CreatedOn: &now,
	})
	return err
}
