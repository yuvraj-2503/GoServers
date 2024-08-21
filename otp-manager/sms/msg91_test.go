package sms

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestMSG91Sender_Send(t *testing.T) {
	type fields struct {
		baseUrl    string
		authKey    string
		templateId string
	}
	type args struct {
		ctx *context.Context
		sms *Sms
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Mobile Otp Test",
			fields: fields{
				baseUrl:    "https://control.msg91.com/api/v5/flow",
				authKey:    "409278AEa5Q6FhOTi65733bdaP1",
				templateId: "657015acd6fc0526436e8f82",
			},
			args: args{
				ctx: nil,
				sms: &Sms{
					Mobile: "+917250378940",
					Otp:    123456,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MSG91Sender{
				baseUrl:    tt.fields.baseUrl,
				authKey:    tt.fields.authKey,
				templateId: tt.fields.templateId,
			}
			if err := m.Send(tt.args.ctx, tt.args.sms); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMSG91Sender_getPayload(t *testing.T) {
	type fields struct {
		baseUrl    string
		authKey    string
		templateId string
	}
	type args struct {
		sms *Sms
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MSG91Sender{
				baseUrl:    tt.fields.baseUrl,
				authKey:    tt.fields.authKey,
				templateId: tt.fields.templateId,
			}
			if got := m.getPayload(tt.args.sms); got != tt.want {
				t.Errorf("getPayload() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMSG91Sender_setHeaders(t *testing.T) {
	type fields struct {
		baseUrl    string
		authKey    string
		templateId string
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MSG91Sender{
				baseUrl:    tt.fields.baseUrl,
				authKey:    tt.fields.authKey,
				templateId: tt.fields.templateId,
			}
			m.setHeaders(tt.args.req)
		})
	}
}

func TestNewMSG91Sender(t *testing.T) {
	type args struct {
		baseUrl    string
		authKey    string
		templateId string
	}
	tests := []struct {
		name string
		args args
		want *MSG91Sender
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMSG91Sender(tt.args.baseUrl, tt.args.authKey, tt.args.templateId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMSG91Sender() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getRecipients(t *testing.T) {
	type args struct {
		sms *Sms
	}
	tests := []struct {
		name string
		args args
		want []Recipient
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRecipients(tt.args.sms); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRecipients() = %v, want %v", got, tt.want)
			}
		})
	}
}
