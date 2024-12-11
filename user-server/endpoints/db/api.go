package db

import (
	"context"
	"time"
)

type UrlData struct {
	Key       string     `bson:"key"`
	Url       string     `bson:"url"`
	Env       string     `bson:"env"`
	UpdatedAt *time.Time `bson:"updated_at"`
}

type UrlStore interface {
	Upsert(ctx *context.Context, urls []*UrlData) error
	Get(ctx *context.Context, key, env string) (*UrlData, error)
	GetAll(ctx *context.Context, env string) ([]*UrlData, error)
	Delete(ctx *context.Context, key, env string) error
}
