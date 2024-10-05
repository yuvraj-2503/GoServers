package api

import "user-server/common"

type SignInRequest struct {
	UserId string         `json:"userId" binding:"required"`
	App    string         `json:"app" binding:"required"`
	Device *common.Device `json:"device" binding:"required"`
}
