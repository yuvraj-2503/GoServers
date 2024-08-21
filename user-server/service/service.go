package service

import (
	"context"
	"user-server/common"
	"user-server/db"
)

type SignUpManager interface {
	SendEmailOtp(ctx *context.Context, emailId string) error
}

type MongoSignupManager struct {
	userStore db.UserStore
}

func NewMongoSignupManager(userStore db.UserStore) *MongoSignupManager {
	return &MongoSignupManager{
		userStore: userStore,
	}
}

func (m *MongoSignupManager) SendEmailOtp(ctx *context.Context, emailId string) error {
	result, _ := m.userStore.CheckExists(ctx, db.Filter{
		Key:   db.EmailId,
		Value: emailId,
	})

	if result {
		return &common.AlreadyExistsError{Message: "Email already registered"}
	}
	return nil
}
