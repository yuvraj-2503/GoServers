package db

import (
    "context"
    "testing"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func Test_createBsonFilter_and_Delete_variants(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping Mongo-backed test in short mode")
    }

    ctx := context.Background()
    clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(ctx, clientOpts)
    if err != nil {
        t.Fatalf("connect err: %v", err)
    }
    defer client.Disconnect(ctx)

    coll := client.Database("user-server-test").Collection("users-extra-tests")
    _ = coll.Drop(ctx)
    defer coll.Drop(ctx)

    store := NewMongoUserStore(coll)

    // create sample user
    user := &User{UserId: "extra-1", EmailId: "e1@example.com"}
    if err := store.Insert(&ctx, user); err != nil {
        t.Fatalf("Insert failed: %v", err)
    }

    // createBsonFilter simple call
    _ = createBsonFilter(Filter{Key: UserId, Value: "extra-1"})

    // Delete by userId using Delete (filter)
    if err := store.Delete(&ctx, Filter{Key: UserId, Value: "extra-1"}); err != nil {
        t.Fatalf("Delete by userId failed: %v", err)
    }

    // Try Delete again to hit NotFound branch
    if err := store.Delete(&ctx, Filter{Key: UserId, Value: "extra-1"}); err == nil {
        t.Fatalf("expected NotFoundError on repeated delete")
    }

    // small pause
    time.Sleep(5 * time.Millisecond)
}
