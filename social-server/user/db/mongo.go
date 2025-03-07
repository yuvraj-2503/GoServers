package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"user-server/common"
)

type MongoUserStore struct {
	collection *mongo.Collection
}

func NewMongoUserStore(collection *mongo.Collection) *MongoUserStore {
	return &MongoUserStore{collection: collection}
}

func (s *MongoUserStore) InsertUser(ctx *context.Context, user *UserDetails) error {
	_, err := s.collection.InsertOne(*ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return &common.AlreadyExistsError{Message: "User already exists"}
		}
		return err
	}
	return nil
}

func (s *MongoUserStore) UpdateUser(ctx *context.Context, user *UserDetails) error {
	filter := bson.D{{
		Key: "userId", Value: user.UserId,
	}}

	result, err := s.collection.UpdateOne(*ctx, filter, user)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return &common.NotFoundError{Message: "User not found"}
	}
	return nil
}

func (s *MongoUserStore) DeleteUser(ctx *context.Context, userId string) error {
	result, err := s.collection.DeleteOne(*ctx, bson.M{"userId": userId})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return &common.NotFoundError{Message: "User not found"}
	}
	return nil
}

func (s *MongoUserStore) GetUserById(ctx *context.Context, userId string) (*UserDetails, error) {
	result := s.collection.FindOne(*ctx, bson.M{"userId": userId})
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, &common.NotFoundError{Message: "User not found"}
	}

	var userDetails UserDetails
	err := result.Decode(&userDetails)
	if err != nil {
		return nil, err
	}
	return &userDetails, nil
}

func (s *MongoUserStore) UpdateFollowersCount(ctx *context.Context, userId string, followersCount int) error {
	filter := bson.D{{
		Key: "userId", Value: userId,
	}}

	update := bson.D{{
		Key: "$inc", Value: bson.D{{
			Key: "followersCount", Value: followersCount,
		}},
	}}

	result, err := s.collection.UpdateOne(*ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return &common.NotFoundError{Message: "User not found"}
	}

	return nil
}

func (s *MongoUserStore) UpdateFollowingCount(ctx *context.Context, userId string, followingCount int) error {
	filter := bson.D{{
		Key: "userId", Value: userId,
	}}
	update := bson.D{{
		Key: "$inc", Value: bson.D{{
			Key: "followingCount", Value: followingCount,
		}},
	}}

	result, err := s.collection.UpdateOne(*ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return &common.NotFoundError{Message: "User not found"}
	}
	return nil
}
