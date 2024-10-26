package mail

import (
	"context"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
)

type SendgridMailSender struct {
	SenderId       string
	SendgridClient *sendgrid.Client
}

func NewSendgridMailSender(senderId, apiKey string) *SendgridMailSender {
	return &SendgridMailSender{
		SenderId:       senderId,
		SendgridClient: sendgrid.NewSendClient(apiKey),
	}
}

func (s *SendgridMailSender) Send(ctx *context.Context, m *TemplatedMail) error {
	message := mail.NewV3Mail()
	message.SetFrom(mail.NewEmail("Yuvraj Singh", s.SenderId))
	p := mail.NewPersonalization()
	setToRecipients(p, m.To)
	setCcRecipients(p, m.Cc)
	setBccRecipients(p, m.Bcc)
	setSendgridTemplate(message, p, m.Template)
	message.AddPersonalizations(p)
	result, err := s.SendgridClient.Send(message)
	log.Println(result)
	if err != nil {
		log.Printf("Failed to send email, reason: %s", err)
		return err
	}
	return nil
}

func setSendgridTemplate(message *mail.SGMailV3, p *mail.Personalization, template *SendgridTemplate) {
	message.SetTemplateID(template.TemplateId)
	for key, value := range template.TemplateData {
		p.SetDynamicTemplateData(key, value)
	}
}

func setToRecipients(p *mail.Personalization, to *[]string) {
	for _, recipient := range *to {
		p.AddTos(mail.NewEmail("", recipient))
	}
}

func setCcRecipients(p *mail.Personalization, cc *[]string) {
	for _, recipient := range *cc {
		p.AddCCs(mail.NewEmail("", recipient))
	}
}

func setBccRecipients(p *mail.Personalization, bcc *[]string) {
	for _, recipient := range *bcc {
		p.AddBCCs(mail.NewEmail("", recipient))
	}
}
