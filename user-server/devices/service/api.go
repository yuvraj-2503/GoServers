package service

import (
	"context"
	"time"
	"user-server/common"
	"user-server/devices/db"
)

type Device struct {
	*common.Device
	UpdatedAt time.Time
}

type UserDeviceManager struct {
	userDeviceStore db.UserDeviceStore
}

func NewUserDeviceManager(userDeviceStore db.UserDeviceStore) *UserDeviceManager {
	return &UserDeviceManager{
		userDeviceStore: userDeviceStore,
	}
}

func (u *UserDeviceManager) RegisterDevice(ctx *context.Context, userId string, device *common.Device) (bool, error) {
	now := time.Now()
	return u.userDeviceStore.Upsert(ctx, &db.UserDevice{
		UserId:     userId,
		DeviceInfo: device,
		UpdatedOn:  &now,
	})
}

func (u *UserDeviceManager) GetDevices(ctx *context.Context, userId string) ([]*Device, error) {
	devices, err := u.userDeviceStore.GetByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	result := []*Device{}
	for _, dev := range devices {
		result = append(result, mapDevice(dev))
	}

	return result, nil
}

func mapDevice(dev *db.UserDevice) *Device {
	return &Device{
		Device:    dev.DeviceInfo,
		UpdatedAt: *dev.UpdatedOn,
	}
}
