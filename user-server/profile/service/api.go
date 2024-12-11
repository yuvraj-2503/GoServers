package service

import (
	blobmanager "blob-manager"
	"bytes"
	"context"
	"encoding/base64"
	"google.golang.org/appengine/log"
	"io"
	"io/ioutil"
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
	GetProfile(ctx *context.Context, userId string, currentTime time.Time) (*UserProfile, error)
	DeleteProfileByUserId(ctx *context.Context, userId string) error
}

type ProfileServiceImpl struct {
	profileStore db.ProfileStore
	blobManager  *blobmanager.BlobManager
}

func NewProfileService(profileStore db.ProfileStore, blobManager *blobmanager.BlobManager) *ProfileServiceImpl {
	return &ProfileServiceImpl{profileStore: profileStore, blobManager: blobManager}
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
	if len(profile.ProfilePicture) > 0 {
		dbProfile.PictureUpdatedOn = &currTime
	}

	var errChan = make(chan error, 2)

	s.updateProfile(ctx, dbProfile, errChan)
	s.updateProfilePic(ctx, userId, profile.ProfilePicture, errChan)
	for i := 0; i < 2; i++ {
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

	return s.mapProfile(ctx, profile)
}

func (s *ProfileServiceImpl) GetProfile(ctx *context.Context, userId string,
	currentTime time.Time) (*UserProfile, error) {
	profile, err := s.profileStore.GetByUserId(ctx, userId, currentTime)
	if err != nil {
		return nil, err
	}

	userProfile := &UserProfile{}
	if profile.UpdatedOn != nil && profile.UpdatedOn.After(currentTime) {
		userProfile.FirstName = profile.FirstName
		userProfile.LastName = profile.LastName
	}

	if profile.PictureUpdatedOn != nil && profile.PictureUpdatedOn.After(currentTime) {
		err := s.setProfilePicture(ctx, userId, userProfile)
		if err != nil {
			return nil, err
		}
	}

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

func (s *ProfileServiceImpl) updateProfilePic(ctx *context.Context,
	userId string, profilePicture string, errChan chan error) {
	go func() {
		if len(profilePicture) == 0 {
			err := s.blobManager.Delete(ctx, userId)
			if err != nil {
				log.Errorf(*ctx, "ERROR: Failed to remove profile picture, reason: %s",
					err.Error())
			}
			errChan <- nil
		}
		if len(profilePicture) > 0 {
			file := getFile(userId, profilePicture)
			err := s.blobManager.Upload(ctx, file)
			if err != nil {
				errChan <- err
			}
			errChan <- nil
		}
	}()
}

func (s *ProfileServiceImpl) setProfilePicture(ctx *context.Context,
	userId string, profile *UserProfile) error {
	file, err := s.downloadFile(ctx, userId)
	if err != nil {
		return err
	}
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	profile.ProfilePicture = base64.StdEncoding.EncodeToString(fileContent)
	return nil
}

func (s *ProfileServiceImpl) downloadFile(ctx *context.Context,
	userId string) (blobmanager.File, error) {
	file, err := s.blobManager.Download(ctx, userId)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s *ProfileServiceImpl) mapProfile(ctx *context.Context, profile *db.Profile) (*UserProfile, error) {
	userProfile := &UserProfile{}
	userProfile.FirstName = profile.FirstName
	userProfile.LastName = profile.LastName

	err := s.setProfilePicture(ctx, profile.UserId, userProfile)
	if err != nil {
		return nil, err
	}
	return userProfile, nil
}

func getFile(userId, fileContent string) blobmanager.File {
	decodedBytes, _ := base64.StdEncoding.DecodeString(fileContent)
	return NewUploadableFile(userId, getReadCloserFromByteArray(decodedBytes), int64(len(decodedBytes)), "image/*")
}

func getReadCloserFromByteArray(data []byte) io.ReadCloser {
	reader := bytes.NewReader(data)
	readCloser := io.NopCloser(reader)
	return readCloser
}
