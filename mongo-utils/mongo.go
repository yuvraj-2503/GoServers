package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoConfig struct {
	ConnectionString string
	Database         string
	Username         string
	Password         string
}

func (c *MongoConfig) getMongoClient() (*mongo.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.ConnectionString))
	return client, err
}

func (c *MongoConfig) GetCollection(collection string) (*mongo.Collection, error) {
	mongoClient, err := c.getMongoClient()
	if err != nil {
		return nil, err
	}
	return mongoClient.Database(c.Database).Collection(collection), nil
}
