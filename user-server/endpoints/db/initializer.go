package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"user-server/common"
	"user-server/config"
)

func LoadDB(ctx *context.Context) {
	mongoConfig := config.Configuration.MongoConfig
	keyEnvIdx := mongo.IndexModel{
		Keys: bson.D{
			{"key", 1},
			{"env", 2},
		},
		Options: options.Index().SetName("key-env-idx").SetUnique(true),
	}
	urlColl, err := mongoConfig.GetCollection(common.UrlsCollection)
	if err != nil {
		log.Panicf("failed to get collection %s , because of %s", common.UrlsCollection, err.Error())
	}
	_, err = urlColl.Indexes().CreateMany(*ctx, []mongo.IndexModel{keyEnvIdx})
	if err != nil {
		log.Panicf("failed to create index on %s, reason: %s", common.UrlsCollection, err)
	}
}
