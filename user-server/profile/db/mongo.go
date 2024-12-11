package db

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"user-server/common"
)

type MongoProfileStore struct {
	profileColl *mongo.Collection
}

func NewMongoProfileStore(profileColl *mongo.Collection) *MongoProfileStore {
	return &MongoProfileStore{
		profileColl: profileColl,
	}
}

func (p *MongoProfileStore) Upsert(ctx *context.Context, profile *Profile) error {
	filter := bson.D{{
		Key: "userId", Value: profile.UserId,
	}}
	updateOpts := options.Update().SetUpsert(true)
	_, err := p.profileColl.UpdateOne(*ctx, filter, getUpdates(profile), updateOpts)
	return err
}

func getUpdates(profile *Profile) bson.D {
	updates := bson.D{}
	if len(profile.FirstName) > 0 {
		updates = append(updates, bson.E{Key: "firstName", Value: profile.FirstName})
	}
	if len(profile.LastName) > 0 {
		updates = append(updates, bson.E{Key: "lastName", Value: profile.LastName})
	}
	if profile.UpdatedOn != nil {
		updates = append(updates, bson.E{Key: "updatedOn", Value: profile.UpdatedOn})
	}
	if profile.PictureUpdatedOn != nil {
		updates = append(updates, bson.E{Key: "pictureUpdatedOn", Value: profile.PictureUpdatedOn})
	}
	return bson.D{bson.E{Key: "$set", Value: updates}}
}

func (p *MongoProfileStore) Get(ctx *context.Context, userId string) (*Profile, error) {
	filter := bson.D{{
		Key: "userId", Value: userId,
	}}
	var profile Profile
	err := p.profileColl.FindOne(*ctx, filter).Decode(&profile)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &common.NotFoundError{Message: "User Profile not found"}
		}
		return nil, err
	}
	return &profile, nil
}

func (p *MongoProfileStore) GetByUserId(ctx *context.Context, userId string,
	time time.Time) (*Profile, error) {
	filter := bson.M{
		"userId": userId,
		"$or": []bson.M{
			{"updatedOn": bson.M{"$gte": time}},
			{"pictureUpdatedOn": bson.M{"$gte": time}},
		},
	}

	var result Profile
	err := p.profileColl.FindOne(*ctx, filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, &common.NotFoundError{
			Message: "User Profile not found",
		}
	}
	return &result, err
}

func (p *MongoProfileStore) Delete(ctx *context.Context, userId string) error {
	filter := bson.D{{
		Key: "userId", Value: userId,
	}}
	result, err := p.profileColl.DeleteOne(*ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return &common.NotFoundError{Message: "User Profile not found"}
	}
	return nil
}
