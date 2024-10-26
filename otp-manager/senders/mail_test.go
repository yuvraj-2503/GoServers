package senders

import (
	"context"
	"mail-sender/mail"
	"otp-manager/common"
	"testing"
)

var ctx context.Context
var sendgridSender mail.TemplatedMailSender

func init() {
	ctx = context.TODO()
	sendgridSender = mail.NewSendgridMailSender("noreply@ziroh.com",
		"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")
}

func TestMailOtpSender_Send(t *testing.T) {
	type fields struct {
		ctx            *context.Context
		sendgridSender mail.TemplatedMailSender
	}
	type args struct {
		contact *common.Contact
		otp     uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Otp Sender Test",
			fields: fields{
				ctx:            &ctx,
				sendgridSender: sendgridSender,
			},
			args: args{
				contact: &common.Contact{
					EmailId: "singh.yuvraj1047@gmail.com",
				},
				otp: 123456,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MailOtpSender{
				ctx:            tt.fields.ctx,
				sendgridSender: tt.fields.sendgridSender,
			}
			if err := s.Send(tt.args.contact, tt.args.otp); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
