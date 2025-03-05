package service

import (
	"context"
	"social-server/common"
	"social-server/user/db"
)

type UserServiceImpl struct {
	db db.UserDB
}

func NewUserService(db db.UserDB) *UserServiceImpl {
	return &UserServiceImpl{db: db}
}

func (s *UserServiceImpl) RegisterUser(ctx *context.Context, details *db.UserDetails) *common.Result {
	err := s.db.InsertUser(ctx, details)
	if err != nil {
		return common.NewResult(1, err.Error())
	}
	return common.NewResult(0, "User Registered Successfully.")
}

func (s *UserServiceImpl) GetUserById(ctx *context.Context, userId string) (*db.UserDetails, error) {
	return s.db.GetUserById(ctx, userId)
}

func (s *UserServiceImpl) UpdateFollowersCount(ctx *context.Context, userId string) error {
	return s.db.UpdateFollowersCount(ctx, userId, 1)
}

func (s *UserServiceImpl) UpdateFollowingCount(ctx *context.Context, userId string) error {
	return s.db.UpdateFollowingCount(ctx, userId, 1)
}
