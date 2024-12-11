package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"user-server/common"
)

type UrlMongoStore struct {
	urlColl *mongo.Collection
}

func NewUrlMongoStore(urlColl *mongo.Collection) *UrlMongoStore {
	return &UrlMongoStore{urlColl: urlColl}
}

func (u *UrlMongoStore) Upsert(ctx *context.Context, urls []*UrlData) error {
	var operations []mongo.WriteModel

	for _, url := range urls {
		filter := bson.D{
			{Key: "key", Value: url.Key},
			{Key: "env", Value: url.Env}}
		update := getUpdates(url)

		// Create an upsert model
		upsertModel := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)

		operations = append(operations, upsertModel)
	}

	// Execute the bulk write
	_, err := u.urlColl.BulkWrite(*ctx, operations, options.BulkWrite().SetOrdered(false))
	return err
}

func (u *UrlMongoStore) Get(ctx *context.Context, key, env string) (*UrlData, error) {
	filter := bson.D{
		{Key: "key", Value: key},
		{Key: "env", Value: env}}

	var result UrlData
	err := u.urlColl.FindOne(*ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &common.NotFoundError{Message: "Url not found"}
		}
		return nil, err
	}
	return &result, nil
}

func (u *UrlMongoStore) GetAll(ctx *context.Context, env string) ([]*UrlData, error) {
	filter := bson.D{
		{Key: "env", Value: env}}
	cursor, err := u.urlColl.Find(*ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &common.NotFoundError{Message: "No Url Data found"}
		}
		return nil, err
	}

	var urlDatas []*UrlData
	defer cursor.Close(*ctx)
	for cursor.Next(*ctx) {
		var urlData UrlData
		err := cursor.Decode(&urlData)
		if err != nil {
			return nil, err
		}

		urlDatas = append(urlDatas, &urlData)
	}
	return urlDatas, nil
}

func (u *UrlMongoStore) Delete(ctx *context.Context, key, env string) error {
	filter := bson.D{
		{Key: "key", Value: key},
		{Key: "env", Value: env}}
	_, err := u.urlColl.DeleteOne(*ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &common.NotFoundError{Message: "Url not found"}
		}
		return err
	}
	return nil
}

func getUpdates(url *UrlData) bson.D {
	updates := bson.D{}
	if len(url.Url) > 0 {
		updates = append(updates, bson.E{Key: "url", Value: url.Url})
	}
	if len(url.Env) > 0 {
		updates = append(updates, bson.E{Key: "env", Value: url.Env})
	}
	if url.UpdatedAt != nil {
		updates = append(updates, bson.E{Key: "updated_at", Value: url.UpdatedAt})
	}
	return bson.D{{Key: "$set", Value: updates}}
}
