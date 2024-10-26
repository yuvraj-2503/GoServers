package mail

import "context"

type TemplatedMailSender interface {
	Send(ctx *context.Context, mail *TemplatedMail) error
}

type SendgridTemplate struct {
	TemplateId   string
	TemplateData map[string]string
}

type TemplatedMail struct {
	To       *[]string
	Cc       *[]string
	Bcc      *[]string
	Template *SendgridTemplate
}
