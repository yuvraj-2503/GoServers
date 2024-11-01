package sessiondb

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"user-server/common"
)

type SessionMongoStore struct {
	collection *mongo.Collection
}

func NewSessionMongoStore(collection *mongo.Collection) *SessionMongoStore {
	return &SessionMongoStore{
		collection: collection,
	}
}

func (s *SessionMongoStore) Insert(ctx *context.Context, sessionMapping *SessionMapping) error {
	_, err := s.collection.InsertOne(*ctx, sessionMapping)
	return err
}

func (s *SessionMongoStore) Get(ctx *context.Context, sessionId string) (*SessionMapping, error) {
	filter := bson.D{{Key: "sessionId", Value: sessionId}}

	var result SessionMapping
	err := s.collection.FindOne(*ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, &common.NotFoundError{}
		}
		return nil, err
	}
	return &result, nil
}

func (s *SessionMongoStore) Delete(ctx *context.Context, sessionId string) error {
	filter := bson.D{{Key: "sessionId", Value: sessionId}}

	result, err := s.collection.DeleteOne(*ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return &common.NotFoundError{}
	}
	return nil
}
