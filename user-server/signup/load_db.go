package signup

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	mongodb "mongo-utils"
	"user-server/common"
	"user-server/config"
)

func LoadDB(ctx *context.Context) {
	mongoConfig := config.Configuration.MongoConfig
	loadOtpCollection(ctx, mongoConfig)
	loadUserCollection(ctx, mongoConfig)
}

func loadUserCollection(ctx *context.Context, mongoConfig *mongodb.MongoConfig) {
	userIdIndex := mongo.IndexModel{
		Keys: bson.D{
			{"userId", 1},
		},
		Options: options.Index().SetName("userId-index").SetUnique(true),
	}
	emailPhoneCompoundIndex := mongo.IndexModel{
		Keys: bson.D{
			{"emailId", 1},
			{"phoneNumber.countryCode", 1},
			{"phoneNumber.number", 1},
		},
		Options: options.Index().SetName("emailId-phone-index").SetUnique(true),
	}
	userColl, _ := mongoConfig.GetCollection(common.UserCollection)
	_, err := userColl.Indexes().CreateMany(*ctx, []mongo.IndexModel{userIdIndex, emailPhoneCompoundIndex})
	if err != nil {
		log.Panicf("failed to create index, reason: %s", err)
	}
}

func loadOtpCollection(ctx *context.Context, mongoConfig *mongodb.MongoConfig) {
	otpCollection, _ := mongoConfig.GetCollection(common.OtpCollection)
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"sessionId", 1},
		},
		Options: options.Index().SetName("sessionId-index").SetUnique(true),
	}
	_, err := otpCollection.Indexes().CreateOne(*ctx, indexModel)
	if err != nil {
		log.Panicf("failed to create index for otp collection, reason: %s", err)
	}
}
