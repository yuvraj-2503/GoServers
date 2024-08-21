package common

import (
	"time"
	"user-server/common"
)

type Contact struct {
	EmailId     string              `json:"emailId" bson:"emailId"`
	PhoneNumber *common.PhoneNumber `json:"phoneNumber" bson:"phoneNumber"`
}

type OTP struct {
	SessionId string    `bson:"sessionId"`
	Contact   *Contact  `bson:"contact"`
	Otp       []byte    `bson:"otp"`
	Retries   uint8     `bson:"retries"`
	CreatedOn time.Time `bson:"createdOn"`
}
