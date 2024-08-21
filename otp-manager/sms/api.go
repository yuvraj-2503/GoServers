package sms

import "context"

type Sms struct {
	Mobile string
	Otp    uint64
}

type SmsSender interface {
	Send(ctx *context.Context, sms *Sms) error
}
