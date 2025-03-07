package service

import (
	"context"
	"time"
	token "token-manager"
	"user-server/common"
	"user-server/signin/api"
	"user-server/signup/db"
)

type SignInManager interface {
	SignIn(ctx *context.Context, request *api.SignInRequest) (*common.AuthenticatedUser, error)
}

type MongoSignInManager struct {
	userStore    db.UserStore
	tokenManager token.TokenManager
}

func NewMongoSignInManager(userStore db.UserStore, tokenManager token.TokenManager) *MongoSignInManager {
	return &MongoSignInManager{userStore: userStore, tokenManager: tokenManager}
}

func (m *MongoSignInManager) SignIn(ctx *context.Context, request *api.SignInRequest) (*common.AuthenticatedUser, error) {
	user, err := m.userStore.Get(ctx, db.Filter{Key: db.UserId, Value: request.UserId})
	if err != nil {
		return nil, &common.UserDoesNotExistError{Message: "User does not exist"}
	}

	authenticatedUser := new(common.AuthenticatedUser)
	authenticatedUser.UserId = user.UserId
	authenticatedUser.EmailId = user.EmailId
	authenticatedUser.PhoneNumber = user.PhoneNumber

	errCh := make(chan error, 1)
	m.createApiKey(authenticatedUser, user.UserId, user.EmailId, request.Device.FingerPrint, request.App, errCh)

	for i := 0; i < 1; i++ {
		e := <-errCh
		if e != nil {
			return nil, e
		}
	}
	return authenticatedUser, nil
}

func (m *MongoSignInManager) createApiKey(
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
		Sub:       userId,
	}
}
