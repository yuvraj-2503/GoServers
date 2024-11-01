package authenticator

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
	sessionColl, _ := mongoConfig.GetCollection(common.SessionUserIdCollection)
	sessionCollIndex := mongo.IndexModel{
		Keys: bson.D{
			{"sessionId", 1},
		},
		Options: options.Index().SetName("Session-Coll-Index").SetUnique(true),
	}

	result, err := sessionColl.Indexes().CreateOne(*ctx, sessionCollIndex)
	if err != nil {
		log.Panicf("Failed to create index on %v, reason %v", common.SessionUserIdCollection, err.Error())
	}
	log.Printf("Crwated index on %s: %s", common.SessionUserIdCollection, result)
}
