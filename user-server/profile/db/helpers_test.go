package db

import (
    "context"
    "testing"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func Test_getUpdates_profile_behaviour_and_GetByUserId_negative(t *testing.T) {
    now := time.Now()
    p := &Profile{UserId: "x", FirstName: "F", UpdatedOn: &now}
    upd := getUpdates(p)
    if len(upd) == 0 {
        t.Fatalf("expected updates non-empty")
    }

    // Negative GetByUserId when not present (real collection)
    ctx := context.Background()
    clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(ctx, clientOpts)
    if err != nil {
        t.Fatalf("failed to connect: %v", err)
    }
    defer client.Disconnect(ctx)

    coll := client.Database("user-server-test").Collection("profiles-negative-test")
    defer coll.Drop(ctx)

    store := NewMongoProfileStore(coll)
    _, err = store.GetByUserId(&ctx, "non-existent-id", time.Now())
    if err == nil {
        t.Fatalf("expected NotFoundError for non-existent id")
    }
}
