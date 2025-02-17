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
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Create MongoDB client options
	clientOptions := options.Client().ApplyURI(c.ConnectionString)

	// If username and password are provided, set authentication credentials
	if c.Username != "" && c.Password != "" {
		credential := options.Credential{
			Username: c.Username,
			Password: c.Password,
		}
		clientOptions.SetAuth(credential)
	}

	client, err := mongo.Connect(ctx, clientOptions)
	return client, err
}

func (c *MongoConfig) GetCollection(collection string) (*mongo.Collection, error) {
	mongoClient, err := c.getMongoClient()
	if err != nil {
		return nil, err
	}
	return mongoClient.Database(c.Database).Collection(collection), nil
}
