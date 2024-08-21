package senders

import (
	"context"
	"otp-manager/common"
	"otp-manager/sms"
)

type SmsOtpSender struct {
	ctx       *context.Context
	smsSender sms.SmsSender
}

func NewSmsOtpSender(ctx *context.Context, smsSender sms.SmsSender) *SmsOtpSender {
	return &SmsOtpSender{ctx: ctx, smsSender: smsSender}
}

func (s *SmsOtpSender) Send(contact *common.Contact, otp uint64) error {
	sms := createSms(contact.PhoneNumber.CountryCode+contact.PhoneNumber.Number, otp)
	return s.smsSender.Send(s.ctx, sms)
}

func createSms(mobile string, otp uint64) *sms.Sms {
	return &sms.Sms{
		Mobile: mobile,
		Otp:    otp,
	}
}
