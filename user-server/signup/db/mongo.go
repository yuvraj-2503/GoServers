package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"user-server/common"
)

type MongoUserStore struct {
	userColl *mongo.Collection
}

func NewMongoUserStore(collection *mongo.Collection) *MongoUserStore {
	return &MongoUserStore{userColl: collection}
}

func (u *MongoUserStore) Insert(ctx *context.Context, user *User) error {
	_, err := u.userColl.InsertOne(*ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return &common.AlreadyExistsError{Message: "User already exists"}
		}
		return err
	}
	return nil
}

func (u *MongoUserStore) Get(ctx *context.Context, filter Filter) (*User, error) {
	bsonFilter := createBsonFilter(filter)
	var user User
	err := u.userColl.FindOne(*ctx, bsonFilter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &common.NotFoundError{Message: "User not found"}
		}
		return nil, err
	}
	return &user, nil
}

func (u *MongoUserStore) GetByPhoneNumber(ctx *context.Context, phoneNumber *common.PhoneNumber) (*User, error) {
	bsonFilter := bson.D{{
		Key: "phoneNumber.countryCode", Value: phoneNumber.CountryCode,
	}, {
		Key: "phoneNumber.number", Value: phoneNumber.Number,
	}}
	var user User
	err := u.userColl.FindOne(*ctx, bsonFilter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &common.NotFoundError{Message: "User not found"}
		}
		return nil, err
	}
	return &user, nil
}

func (u *MongoUserStore) UpdateEmailId(ctx *context.Context, userId, emailId string) error {
	filter := bson.D{{Key: "userId", Value: userId}}
	updates := bson.D{{Key: "$set", Value: bson.D{{Key: "emailId", Value: emailId}}}}
	_, err := u.userColl.UpdateOne(*ctx, filter, updates)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return &common.AlreadyExistsError{Message: "User already exists"}
		}
		return err
	}
	return nil
}

func (u *MongoUserStore) UpdatePhoneNumber(ctx *context.Context, userId string, phoneNumber *common.PhoneNumber) error {
	filter := bson.D{{Key: "userId", Value: userId}}
	updates := bson.D{{Key: "$set", Value: bson.D{{Key: "phoneNumber", Value: phoneNumber}}}}
	_, err := u.userColl.UpdateOne(*ctx, filter, updates)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return &common.AlreadyExistsError{Message: "User already exists"}
		}
		return err
	}
	return nil
}

func (u *MongoUserStore) DeleteByUserId(ctx *context.Context, userId string) error {
	filter := bson.D{{Key: "userId", Value: userId}}
	deleteResult, err := u.userColl.DeleteOne(*ctx, filter)
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount == 0 {
		return &common.NotFoundError{Message: "User not found"}
	}
	return nil
}

func (u *MongoUserStore) Delete(ctx *context.Context, filter Filter) error {
	bsonFilter := createBsonFilter(filter)
	deleteResult, err := u.userColl.DeleteOne(*ctx, bsonFilter)
	if err != nil {
		return err
	}
	if deleteResult.DeletedCount == 0 {
		return &common.NotFoundError{Message: "User not found"}
	}
	return nil
}

func (u *MongoUserStore) CheckExists(ctx *context.Context, filter Filter) (bool, error) {
	bsonFilter := createBsonFilter(filter)
	result, err := u.userColl.CountDocuments(*ctx, bsonFilter)
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (u *MongoUserStore) CheckIfMobileExists(ctx *context.Context, phoneNumber *common.PhoneNumber) (bool, error) {
	bsonFilter := bson.D{{
		Key: "phoneNumber.countryCode", Value: phoneNumber.CountryCode,
	}, {
		Key: "phoneNumber.number", Value: phoneNumber.Number,
	}}
	result, err := u.userColl.CountDocuments(*ctx, bsonFilter)
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func createBsonFilter(filter Filter) bson.D {
	return bson.D{{
		Key:   string(filter.Key),
		Value: filter.Value,
	}}
}
