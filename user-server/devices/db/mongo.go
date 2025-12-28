package db

import (
	"context"
	"user-server/common"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUserDeviceStore struct {
	userDevicesCollection *mongo.Collection
}

func NewMongoUserDeviceStore(collection *mongo.Collection) *MongoUserDeviceStore {
	return &MongoUserDeviceStore{
		userDevicesCollection: collection,
	}
}

func (s *MongoUserDeviceStore) Upsert(ctx *context.Context, device *UserDevice) (bool, error) {
	filter := getDeviceMatchQuery(device.UserId, device.DeviceInfo.FingerPrint)
	updateOptions := options.Update().SetUpsert(true)

	updates := getUpdates(device)
	result, err := s.userDevicesCollection.UpdateOne(*ctx, filter, bson.D{{Key: "$set", Value: updates}}, updateOptions)
	if err != nil {
		return false, err
	}

	return result.UpsertedCount > 0, nil
}

func getUpdates(device *UserDevice) bson.D {
	updates := bson.D{}
	if device.DeviceInfo != nil {
		updates = append(updates, bson.E{Key: "deviceInfo", Value: device.DeviceInfo})
	}

	if device.UpdatedOn != nil {
		updates = append(updates, bson.E{Key: "updatedOn", Value: device.UpdatedOn})
	}

	return updates
}

func (s *MongoUserDeviceStore) GetByUserId(ctx *context.Context, userId string) ([]*UserDevice, error) {
	filterQuery := bson.D{
		{Key: "userId", Value: userId},
	}

	cursor, err := s.userDevicesCollection.Find(*ctx, filterQuery)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(*ctx)

	devices := []*UserDevice{}
	if err := cursor.All(*ctx, &devices); err != nil {
		return nil, err
	}

	return devices, nil
}

func (s *MongoUserDeviceStore) Delete(ctx *context.Context, userId string, fingerPrint string) error {
	filterQuery := getDeviceMatchQuery(userId, fingerPrint)

	deleteCount, err := s.userDevicesCollection.DeleteOne(*ctx, filterQuery)
	if err != nil {
		return err
	}

	if deleteCount.DeletedCount == 0 {
		return &common.NotFoundError{}
	}

	return nil
}

func (s *MongoUserDeviceStore) DeleteAllByUserId(ctx *context.Context, userId string) error {
	filterQuery := bson.D{{Key: "userId", Value: userId}}

	deleteCount, err := s.userDevicesCollection.DeleteMany(*ctx, filterQuery)
	if err != nil {
		return err
	}

	if deleteCount.DeletedCount == 0 {
		return &common.NotFoundError{}
	}

	return nil
}

func getDeviceMatchQuery(userId, fingerPrint string) bson.D {
	return bson.D{
		{
			Key:   "userId",
			Value: userId,
		},
		{
			Key:   "deviceInfo.fingerPrint",
			Value: fingerPrint,
		},
	}
}
