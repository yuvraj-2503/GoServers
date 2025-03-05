package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	mongodb "mongo-utils"
	"social-server/common"
	"social-server/config"
)

func LoadDB(ctx *context.Context) {
	mongoConfig := config.Configuration.MongoConfig
	loadUserCollection(ctx, mongoConfig)
}

func loadUserCollection(ctx *context.Context, mongoConfig *mongodb.MongoConfig) {
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"userId", 1},
		},
		Options: options.Index().SetName("userId-Index").SetUnique(true),
	}

	userColl, err := mongoConfig.GetCollection(common.UsersCollection)
	if err != nil {
		log.Fatalf("Failed to connect to database, reason: %v", err.Error())
	}

	indexes, err := userColl.Indexes().CreateMany(*ctx, []mongo.IndexModel{indexModel})
	if err != nil {
		log.Fatalf("Failed to create indexes, reason: %v", err.Error())
	}
	log.Printf("Successfully created indexes with name: %v", indexes[0])
}
