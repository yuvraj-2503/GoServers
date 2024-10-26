package mail

import (
	"context"
	"github.com/sendgrid/sendgrid-go"
	"strconv"
	"testing"
)

func TestSendgridMailSender_Send(t *testing.T) {
	type fields struct {
		SenderId       string
		SendgridClient *sendgrid.Client
	}
	type args struct {
		ctx *context.Context
		m   *TemplatedMail
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Send Email Test",
			fields: fields{
				SenderId:       "noreply@ziroh.com",
				SendgridClient: sendgrid.NewSendClient("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"),
			},
			args: args{
				ctx: nil,
				m: &TemplatedMail{
					To:  &[]string{"singh.yuvraj.yuvraj09@gmail.com"},
					Cc:  &[]string{},
					Bcc: &[]string{},
					Template: &SendgridTemplate{
						TemplateId: "d-f990ca55eab844f49b9763e15b4cead8",
						TemplateData: map[string]string{
							"otp": strconv.Itoa(123456),
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SendgridMailSender{
				SenderId:       tt.fields.SenderId,
				SendgridClient: tt.fields.SendgridClient,
			}
			if err := s.Send(tt.args.ctx, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
