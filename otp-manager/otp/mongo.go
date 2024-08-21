package otp

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"otp-manager/common"
	otpErrors "otp-manager/errors"
)

type MongoOTPStore struct {
	Collection *mongo.Collection
}

func NewMongoOTPStore(collection *mongo.Collection) *MongoOTPStore {
	return &MongoOTPStore{
		Collection: collection,
	}
}

func (store *MongoOTPStore) Upsert(ctx *context.Context, otp *common.OTP) error {
	filter := bson.D{{
		"sessionId",
		otp.SessionId,
	}}
	update := bson.D{
		{"$set", bson.D{
			{"otp", otp.Otp},
			{"contact", otp.Contact},
			{"retries", otp.Retries},
			{"createdOn", otp.CreatedOn},
		}},
	}

	options := options.Update().SetUpsert(true)
	_, err := store.Collection.UpdateOne(*ctx, filter, update, options)
	return err
}

func (store *MongoOTPStore) Delete(ctx *context.Context, sessionId string) error {
	filter := bson.D{{
		"sessionId", sessionId,
	}}
	deleteResult, err := store.Collection.DeleteOne(*ctx, filter)
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount == 0 {
		return &otpErrors.OtpDoesNotExistError{}
	}
	return err
}

func (store *MongoOTPStore) Get(ctx *context.Context, sessionId string) (*common.OTP, error) {
	filter := bson.D{{
		"sessionId", sessionId,
	}}

	var otp common.OTP
	err := store.Collection.FindOne(*ctx, filter).Decode(&otp)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &otpErrors.OtpDoesNotExistError{}
		}
		return nil, err
	}
	return &otp, nil
}
