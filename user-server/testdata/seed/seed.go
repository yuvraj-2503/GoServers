package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println("mongo connect error:", err)
		return
	}
	defer client.Disconnect(ctx)
	db := client.Database("user-server")

	// seed profile
	profileColl := db.Collection("profile-collection")
	profile := bson.M{
		"userId":    "6719252c58da11805939fea3",
		"firstName": "Yuvraj",
		"lastName":  "Singh Rajpoot",
		"updatedOn": time.Date(2024, time.Month(11), 28, 0, 30, 30, 0, time.UTC),
	}
	_, err = profileColl.UpdateOne(ctx, bson.M{"userId": profile["userId"]}, bson.M{"$set": profile}, options.Update().SetUpsert(true))
	if err != nil {
		fmt.Println("seed profile error:", err)
	} else {
		fmt.Println("seeded profile-collection")
	}

	// seed urls collection
	urlColl := db.Collection("urls")
	url := bson.M{
		"key":       "user-server",
		"url":       "http://localhost:8080/api/v1",
		"env":       "LOCAL",
		"updatedAt": time.Now(),
	}
	_, err = urlColl.UpdateOne(ctx, bson.M{"key": url["key"], "env": url["env"]}, bson.M{"$set": url}, options.Update().SetUpsert(true))
	if err != nil {
		fmt.Println("seed urls error:", err)
	} else {
		fmt.Println("seeded urls collection")
	}
}
