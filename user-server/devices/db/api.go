package db

import (
	"context"
	"time"
	"user-server/common"
)

type UserDevice struct {
	UserId     string         `bson:"user_id"`
	DeviceInfo *common.Device `bson:"device"`
	UpdatedOn  *time.Time     `bson:"updated_on"`
}

type UserDeviceStore interface {
	// Upsert will update the record based on userId & fingeprint id combination. If no document exists matching the userId-fingerprint id combination then insert.
	// It returns the boolean value (true if insert otherwise false) or error if any
	Upsert(ctx *context.Context, device *UserDevice) (bool, error)
	GetByUserId(ctx *context.Context, userId string) ([]*UserDevice, error)
	Delete(ctx *context.Context, userId string, fingerPrint string) error
	DeleteAllByUserId(ctx *context.Context, userId string) error
}
