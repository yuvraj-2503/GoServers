package db

import (
	"context"
	"user-server/common"
)

type UserDetails struct {
	UserId         string              `json:"userId" bson:"userId"`
	UserName       string              `json:"userName" bson:"userName"`
	FirstName      string              `json:"firstName" bson:"firstName"`
	LastName       string              `json:"lastName" bson:"lastName"`
	Email          string              `json:"email" bson:"email"`
	PhoneNumber    *common.PhoneNumber `json:"phoneNumber" bson:"phoneNumber"`
	FollowersCount int                 `json:"followersCount" bson:"followersCount"`
	FollowingCount int                 `json:"followingCount" bson:"followingCount"`
}

type UserDB interface {
	InsertUser(ctx *context.Context, user *UserDetails) error
	UpdateUser(ctx *context.Context, user *UserDetails) error
	DeleteUser(ctx *context.Context, userId string) error
	GetUserById(ctx *context.Context, userId string) (*UserDetails, error)
	UpdateFollowersCount(ctx *context.Context, userId string, followersCount int) error
	UpdateFollowingCount(ctx *context.Context, userId string, followingCount int) error
}
