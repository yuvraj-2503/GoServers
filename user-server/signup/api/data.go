package api

import "user-server/common"

type SignUpRequest struct {
	App         string              `json:"app" binding:"required"`
	EmailId     string              `json:"emailId"`
	PhoneNumber *common.PhoneNumber `json:"phoneNumber"`
	OTP         uint64              `json:"otp" binding:"required"`
	Device      *common.Device      `json:"device" binding:"required"`
	SessionId   string              `json:"sessionId" binding:"required"`
}

type EmailOTPRequest struct {
	EmailId string `json:"emailId" binding:"required"`
}
