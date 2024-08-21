package otp

import (
	"context"
	"log"
	mongodb "mongo-utils"
	"otp-manager/common"
	"otp-manager/senders"
	"otp-manager/sms"
	"reflect"
	"testing"
	user "user-server/common"
)

var otpStore OTPStore
var otpSender senders.OtpSender
var ctx context.Context

func init() {
	ctx = context.Background()
	config := mongodb.MongoConfig{
		ConnectionString: "mongodb://localhost:27017",
		Database:         "otp-manager",
		Username:         "",
		Password:         "",
	}
	coll, err := config.GetCollection("otp")
	if err != nil {
		panic(err)
	}
	otpStore = NewMongoOTPStore(coll)
	smsSender := sms.NewMSG91Sender("https://control.msg91.com/api/v5/flow", "409278AEa5Q6FhOTi65733bdaP1",
		"657015acd6fc0526436e8f82")
	otpSender = senders.NewSmsOtpSender(&ctx, smsSender)
}

func TestMongoOtpManager_Send(t *testing.T) {
	result := ""
	type fields struct {
		otpStore OTPStore
		sender   senders.OtpSender
	}
	type args struct {
		ctx     *context.Context
		contact *common.Contact
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Send OTP",
			fields: fields{
				otpStore: otpStore,
				sender:   otpSender,
			},
			args: args{
				ctx: &ctx,
				contact: &common.Contact{
					PhoneNumber: &user.PhoneNumber{
						CountryCode: "+91",
						Number:      "7250378940",
					},
				},
			},
			want:    &result,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MongoOtpManager{
				otpStore: tt.fields.otpStore,
				sender:   tt.fields.sender,
			}
			got, err := m.Send(tt.args.ctx, tt.args.contact)
			if (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				log.Println("SessionId : " + *got)
				t.Errorf("Send() got = %v, want %v", got, tt.want)
			}
		})
	}
}
