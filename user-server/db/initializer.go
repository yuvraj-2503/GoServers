package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	mongodb "mongo-utils"
)

func LoadDb(ctx context.Context, config *mongodb.MongoConfig) {
	usersIdx := mongo.IndexModel{
		Keys: bson.D{
			{"userId", 1},
		},
		Options: options.Index().SetName("Users_idx").SetUnique(true),
	}

	userColl, err := config.GetCollection("users")
	if err != nil {
		log.Panicf("Error loading users collection: %v", err)
	}

	_, err = userColl.Indexes().CreateMany(ctx, []mongo.IndexModel{usersIdx})
	if err != nil {
		log.Panicf("Error creating users index: %v", err)
	}
}
