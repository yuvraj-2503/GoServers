package senders

import (
	"context"
	"mail-sender/mail"
	"otp-manager/common"
	"strconv"
)

type MailOtpSender struct {
	ctx            *context.Context
	sendgridSender mail.TemplatedMailSender
}

func NewMailOtpSender(ctx *context.Context, sendgridSender mail.TemplatedMailSender) *MailOtpSender {
	return &MailOtpSender{
		ctx:            ctx,
		sendgridSender: sendgridSender,
	}
}

func (s *MailOtpSender) Send(contact *common.Contact, otp uint64) error {
	sendGridMail := s.createMail(contact.EmailId, otp)
	return s.sendgridSender.Send(s.ctx, sendGridMail)
}

func (s *MailOtpSender) createMail(emailId string, otp uint64) *mail.TemplatedMail {
	return &mail.TemplatedMail{
		To:       &[]string{emailId},
		Cc:       &[]string{},
		Bcc:      &[]string{},
		Template: s.createTemplate(otp),
	}
}

func (s *MailOtpSender) createTemplate(otp uint64) *mail.SendgridTemplate {
	return &mail.SendgridTemplate{
		TemplateId:   "d-f990ca55eab844f49b9763e15b4cead8",
		TemplateData: map[string]string{"otp": strconv.Itoa(int(otp))},
	}
}
