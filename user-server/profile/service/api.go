package service

import (
	"context"
	"time"
	"user-server/profile/db"
)

type UserProfile struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	ProfilePicture string `json:"profilePicture"`
}

type ProfileService interface {
	UpsertProfile(ctx *context.Context, userId string, profile *UserProfile) error
	GetProfileByUserId(ctx *context.Context, userId string) (*UserProfile, error)
	DeleteProfileByUserId(ctx *context.Context, userId string) error
}

type ProfileServiceImpl struct {
	profileStore db.ProfileStore
}

func NewProfileService(profileStore db.ProfileStore) *ProfileServiceImpl {
	return &ProfileServiceImpl{profileStore: profileStore}
}

func (s *ProfileServiceImpl) UpsertProfile(ctx *context.Context,
	userId string, profile *UserProfile) error {
	dbProfile := &db.Profile{
		UserId:    userId,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
	}
	var currTime = time.Now()
	if len(profile.FirstName) > 0 || len(profile.LastName) > 0 {
		dbProfile.UpdatedOn = &currTime
	}

	var errChan = make(chan error, 1)

	s.updateProfile(ctx, dbProfile, errChan)
	for i := 0; i < 1; i++ {
		err := <-errChan
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ProfileServiceImpl) GetProfileByUserId(ctx *context.Context,
	userId string) (*UserProfile, error) {
	profile, err := s.profileStore.Get(ctx, userId)
	if err != nil {
		return nil, err
	}

	userProfile := &UserProfile{}
	userProfile.FirstName = profile.FirstName
	userProfile.LastName = profile.LastName
	return userProfile, nil
}

func (s *ProfileServiceImpl) DeleteProfileByUserId(ctx *context.Context, userId string) error {
	return nil
}

func (s *ProfileServiceImpl) updateProfile(ctx *context.Context, profile *db.Profile,
	errChan chan error) {
	go func() {
		errChan <- s.profileStore.Upsert(ctx, profile)
	}()
}
