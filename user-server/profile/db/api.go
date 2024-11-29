package db

import (
	"context"
	"time"
)

type Profile struct {
	UserId    string     `bson:"userId"`
	FirstName string     `bson:"firstName"`
	LastName  string     `bson:"lastName"`
	UpdatedOn *time.Time `bson:"updatedOn"`
}

type ProfileStore interface {
	Upsert(ctx *context.Context, profile *Profile) error
	Get(ctx *context.Context, userId string) (*Profile, error)
	Delete(ctx *context.Context, userId string) error
}
