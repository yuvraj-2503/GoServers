package db

import (
	"context"
	"time"
)

type Profile struct {
	UserId           string     `bson:"userId"`
	FirstName        string     `bson:"firstName"`
	LastName         string     `bson:"lastName"`
	UpdatedOn        *time.Time `bson:"updatedOn"`
	PictureUpdatedOn *time.Time `bson:"pictureUpdatedOn"`
}

type ProfileStore interface {
	Upsert(ctx *context.Context, profile *Profile) error
	Get(ctx *context.Context, userId string) (*Profile, error)
	GetByUserId(ctx *context.Context, userId string, time time.Time) (*Profile, error)
	Delete(ctx *context.Context, userId string) error
}
