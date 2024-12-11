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
	userIdIndex := mongo.IndexModel{
		Keys: bson.D{
			{"userId", 1},
		},
		Options: options.Index().SetName("userId-index").SetUnique(true),
	}
	updatedOnIndex := mongo.IndexModel{
		Keys: bson.D{
			{"updatedOn", 1},
			{"pictureUpdatedOn", 1},
		},
		Options: options.Index().SetName("time-index").SetUnique(false),
	}
	profileColl, err := mongoConfig.GetCollection(common.ProfileCollection)
	if err != nil {
		log.Panicf("failed to get collection %s , because of %s", common.ProfileCollection, err.Error())
	}
	_, err = profileColl.Indexes().CreateMany(*ctx, []mongo.IndexModel{userIdIndex, updatedOnIndex})
	if err != nil {
		log.Panicf("failed to create index on %s, reason: %s", common.ProfileCollection, err)
	}
}
