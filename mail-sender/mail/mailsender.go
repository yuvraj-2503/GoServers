package mail

import (
	"context"
	"fmt"
	"github.com/mailersend/mailersend-go"
)

type Mail struct {
	From    string
	To      []string
	Subject string
	Content string
}

type MailSender interface {
	Send(ctx *context.Context, mail *Mail) error
}

type MailSenderImpl struct {
	apiKey string
}

func NewMailSenderImpl(apiKey string) *MailSenderImpl {
	return &MailSenderImpl{
		apiKey: apiKey,
	}
}

func (m *MailSenderImpl) Send(ctx *context.Context, mail *Mail) error {
	recipients := []mailersend.Recipient{}
	for _, to := range mail.To {
		recipients = append(recipients, mailersend.Recipient{
			Email: to,
		})
	}

	from := mailersend.From{
		Name:  "Vibely",
		Email: mail.From,
	}

	ms := mailersend.NewMailersend(m.apiKey)
	message := ms.Email.NewMessage()
	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(mail.Subject)
	message.SetHTML(mail.Content)

	res, err := ms.Email.Send(*ctx, message)
	if err != nil {
		return err
	}
	fmt.Println(res.Body)
	return nil
}
