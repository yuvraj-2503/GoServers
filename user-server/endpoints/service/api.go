package service

import (
	"context"
	"time"
	"user-server/endpoints/db"
)

type UrlService interface {
	Upsert(ctx *context.Context, urls map[string]string, env string) error
	GetAll(ctx *context.Context, env string) (map[string]string, error)
}

type UrlServiceImpl struct {
	urlStore db.UrlStore
}

func NewUrlService(urlStore db.UrlStore) *UrlServiceImpl {
	return &UrlServiceImpl{urlStore: urlStore}
}

func (u *UrlServiceImpl) Upsert(ctx *context.Context, urls map[string]string, env string) error {
	var updatedAt = time.Now()
	dbUrls := []*db.UrlData{}
	for key, url := range urls {
		dbUrl := &db.UrlData{
			Key:       key,
			Url:       url,
			UpdatedAt: &updatedAt,
			Env:       env,
		}
		dbUrls = append(dbUrls, dbUrl)
	}
	return u.urlStore.Upsert(ctx, dbUrls)
}

func (u *UrlServiceImpl) GetAll(ctx *context.Context, env string) (map[string]string, error) {
	dbUrls, err := u.urlStore.GetAll(ctx, env)
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	for _, dbUrl := range dbUrls {
		result[dbUrl.Key] = dbUrl.Url
	}
	return result, nil
}
