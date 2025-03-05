package service

import (
	"context"
	"social-server/common"
	"social-server/user/db"
)

type UserService interface {
	RegisterUser(ctx *context.Context, user *db.UserDetails) *common.Result
	GetUserById(ctx *context.Context, userId string) (*db.UserDetails, error)
	UpdateFollowersCount(ctx *context.Context, userId string) error
	UpdateFollowingCount(ctx *context.Context, userId string) error
}
