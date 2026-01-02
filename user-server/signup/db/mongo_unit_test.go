package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"user-server/common"
)

func getMongoClientLocal(ctx context.Context) (*mongo.Client, error) {
    clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
    return mongo.Connect(ctx, clientOpts)
}

func TestMongoUserStore_Unit_InsertGetCheckExists_GetByPhoneNumber(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping Mongo-backed unit test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("failed to connect to MongoDB: %v", err)
    }
    defer client.Disconnect(ctx)

    collName := fmt.Sprintf("users-unit-test-%d", time.Now().UnixNano())
    coll := client.Database("user-server-test").Collection(collName)
    defer coll.Drop(ctx)

    store := NewMongoUserStore(coll)

    phone := &common.PhoneNumber{CountryCode: "+1", Number: "5551234567"}
    user := &User{
        UserId:      "unit-user-1",
        EmailId:     "unit@example.com",
        PhoneNumber: phone,
    }

    // Insert
    if err := store.Insert(&ctx, user); err != nil {
        t.Fatalf("Insert() error = %v", err)
    }

    // Get
    got, err := store.Get(&ctx, Filter{Key: UserId, Value: "unit-user-1"})
    if err != nil {
        t.Fatalf("Get() error = %v", err)
    }
    if got == nil || got.UserId != "unit-user-1" {
        t.Fatalf("Get() returned unexpected result: %+v", got)
    }

    // GetByPhoneNumber
    got2, err := store.GetByPhoneNumber(&ctx, phone)
    if err != nil {
        t.Fatalf("GetByPhoneNumber() error = %v", err)
    }
    if got2 == nil || got2.UserId != "unit-user-1" {
        t.Fatalf("GetByPhoneNumber() returned unexpected result: %+v", got2)
    }

    // CheckExists
    exists, err := store.CheckExists(&ctx, Filter{Key: UserId, Value: "unit-user-1"})
    if err != nil {
        t.Fatalf("CheckExists() error = %v", err)
    }
    if !exists {
        t.Fatalf("CheckExists() returned false, expected true")
    }
}
