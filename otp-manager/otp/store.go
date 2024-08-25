package otp

import (
	"context"
	"otp-manager/common"
)

type OTPStore interface {
	Upsert(ctx *context.Context, otp *common.OTP) error
	Delete(ctx *context.Context, sessionId string) error
	Get(ctx *context.Context, sessionId string) (*common.OTP, error)
	UpdateRetryCount(ctx *context.Context, sessionId string) error
}
