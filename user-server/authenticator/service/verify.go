package service

import (
	"context"
	"time"
	token "token-manager"
)

func (manager *UserAuthenticatorImpl) Verify(ctx *context.Context, sessionId string, otp uint64) (*VerifyResponse, error) {
	result, err := manager.sessionMapping.Get(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	_, err = manager.emailOtpManager.Verify(ctx, result.SessionId, otp)
	if err != nil {
		return nil, err
	}

	claims := generateTokenClaims(result.UserId)
	authToken, err := manager.tokenManager.Generate(claims)
	if err != nil {
		return nil, err
	}

	return &VerifyResponse{
		UserId: result.UserId,
		Token:  authToken,
	}, nil
}

func generateTokenClaims(userId string) *token.TokenClaims {
	iat := time.Now()
	exp := iat.Add(24 * 7 * time.Hour)
	return &token.TokenClaims{
		UserId: userId,
		IAT:    &iat,
		EXP:    &exp,
		Kind:   "USER",
		Sub:    userId,
	}
}
