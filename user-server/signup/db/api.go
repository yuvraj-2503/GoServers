package db

import (
	"context"
	"user-server/common"
)

type User struct {
	UserId      string              `json:"userId" bson:"userId"`
	EmailId     string              `json:"emailId" bson:"emailId"`
	PhoneNumber *common.PhoneNumber `json:"phoneNumber" bson:"phoneNumber"`
}

type SearchKey string

const (
	UserId  SearchKey = "userId"
	EmailId SearchKey = "emailId"
)

type Filter struct {
	Key   SearchKey
	Value string
}

type UserStore interface {
	Insert(ctx *context.Context, user *User) error
	Get(ctx *context.Context, filter Filter) (*User, error)
	GetByPhoneNumber(ctx *context.Context, phoneNumber *common.PhoneNumber) (*User, error)
	UpdateEmailId(ctx *context.Context, userId, emailId string) error
	UpdatePhoneNumber(ctx *context.Context, userId string, phoneNumber *common.PhoneNumber) error
	DeleteByUserId(ctx *context.Context, userId string) error
	Delete(ctx *context.Context, filter Filter) error
	CheckExists(ctx *context.Context, filter Filter) (bool, error)
	CheckIfMobileExists(ctx *context.Context, phoneNumber *common.PhoneNumber) (bool, error)
}
