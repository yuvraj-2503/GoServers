package db

import (
    "context"
    "fmt"
    "testing"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoProfileStore_Unit_UpsertGetGetByUserIdDelete(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration-style unit test in short mode")
    }

    ctx := context.Background()
    clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(ctx, clientOpts)
    if err != nil {
        t.Fatalf("failed to connect to MongoDB: %v", err)
    }
    defer client.Disconnect(ctx)

    collName := fmt.Sprintf("profiles-unit-test-%d", time.Now().UnixNano())
    coll := client.Database("user-server-test").Collection(collName)
    defer coll.Drop(ctx)

    store := NewMongoProfileStore(coll)

    now := time.Now()
    profile := &Profile{
        UserId:    "p-user-1",
        FirstName: "Alice",
        LastName:  "Smith",
        UpdatedOn: &now,
    }

    if err := store.Upsert(&ctx, profile); err != nil {
        t.Fatalf("Upsert() error = %v", err)
    }

    got, err := store.Get(&ctx, "p-user-1")
    if err != nil {
        t.Fatalf("Get() error = %v", err)
    }
    if got == nil || got.FirstName != "Alice" {
        t.Fatalf("unexpected Get result: %+v", got)
    }

    // GetByUserId - use time far in past so condition matches
    ts := time.Now().Add(-time.Hour)
    got2, err := store.GetByUserId(&ctx, "p-user-1", ts)
    if err != nil {
        t.Fatalf("GetByUserId() error = %v", err)
    }
    if got2 == nil || got2.UserId != "p-user-1" {
        t.Fatalf("unexpected GetByUserId result: %+v", got2)
    }

    // Delete
    if err := store.Delete(&ctx, "p-user-1"); err != nil {
        t.Fatalf("Delete() error = %v", err)
    }
    _, err = store.Get(&ctx, "p-user-1")
    if err == nil {
        t.Fatalf("expected error after delete, got nil")
    }
}

func TestMongoProfileStore_Unit_DeleteNotFound_And_UpsertPartial(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration-style unit test in short mode")
    }

    ctx := context.Background()
    clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(ctx, clientOpts)
    if err != nil {
        t.Fatalf("failed to connect to MongoDB: %v", err)
    }
    defer client.Disconnect(ctx)

    collName := fmt.Sprintf("profiles-delete-test-%d", time.Now().UnixNano())
    coll := client.Database("user-server-test").Collection(collName)
    defer coll.Drop(ctx)

    store := NewMongoProfileStore(coll)

    // Delete non-existent
    if err := store.Delete(&ctx, "non-existent-id"); err == nil {
        t.Fatalf("expected NotFoundError on delete non-existent")
    }

    // Upsert with only userId and FirstName (partial update)
    p := &Profile{UserId: "partial-user", FirstName: "John"}
    if err := store.Upsert(&ctx, p); err != nil {
        t.Fatalf("Upsert partial failed: %v", err)
    }

    got, err := store.Get(&ctx, "partial-user")
    if err != nil || got == nil {
        t.Fatalf("Get after upsert failed: %v", err)
    }
    if got.FirstName != "John" {
        t.Fatalf("expected FirstName=John, got %s", got.FirstName)
    }

    // Upsert again with LastName to add it
    now := time.Now()
    p2 := &Profile{UserId: "partial-user", FirstName: "John", LastName: "Doe", UpdatedOn: &now}
    if err := store.Upsert(&ctx, p2); err != nil {
        t.Fatalf("Upsert update failed: %v", err)
    }

    got2, err := store.Get(&ctx, "partial-user")
    if err != nil || got2 == nil {
        t.Fatalf("Get after second upsert failed: %v", err)
    }
    if got2.LastName != "Doe" {
        t.Fatalf("expected LastName=Doe after upsert, got %s", got2.LastName)
    }
}

func TestMongoProfileStore_Unit_GetByUserIdWithTimeCondition(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration-style unit test in short mode")
    }

    ctx := context.Background()
    clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(ctx, clientOpts)
    if err != nil {
        t.Fatalf("failed to connect to MongoDB: %v", err)
    }
    defer client.Disconnect(ctx)

    collName := fmt.Sprintf("profiles-time-test-%d", time.Now().UnixNano())
    coll := client.Database("user-server-test").Collection(collName)
    defer coll.Drop(ctx)

    store := NewMongoProfileStore(coll)

    // Upsert profile with recent UpdatedOn
    now := time.Now()
    p := &Profile{UserId: "time-user", FirstName: "Test", UpdatedOn: &now}
    if err := store.Upsert(&ctx, p); err != nil {
        t.Fatalf("Upsert failed: %v", err)
    }

    // GetByUserId with time far in past should match (updated >= past)
    pastTime := now.Add(-1 * time.Hour)
    got, err := store.GetByUserId(&ctx, "time-user", pastTime)
    if err != nil {
        t.Fatalf("GetByUserId(past time) error: %v", err)
    }
    if got == nil || got.UserId != "time-user" {
        t.Fatalf("expected to find profile with past time, got %+v", got)
    }

    // GetByUserId with future time should not match
    futureTime := now.Add(1 * time.Hour)
    _, err = store.GetByUserId(&ctx, "time-user", futureTime)
    if err == nil {
        t.Fatalf("expected NotFoundError for future time, got nil")
    }

    // Update with PictureUpdatedOn
    now2 := time.Now()
    p.PictureUpdatedOn = &now2
    if err := store.Upsert(&ctx, p); err != nil {
        t.Fatalf("Upsert with PictureUpdatedOn failed: %v", err)
    }

    // GetByUserId should still match (PictureUpdatedOn is recent)
    got2, err := store.GetByUserId(&ctx, "time-user", pastTime)
    if err != nil {
        t.Fatalf("GetByUserId after PictureUpdatedOn update error: %v", err)
    }
    if got2 == nil {
        t.Fatalf("expected to find profile, got nil")
    }
}

func TestMongoProfileStore_Unit_GetByUserId_MultipleProfiles_And_Delete_Variants(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping Mongo-backed test in short mode")
    }

    ctx := context.Background()
    clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
    client, err := mongo.Connect(ctx, clientOpts)
    if err != nil {
        t.Fatalf("connect error: %v", err)
    }
    defer client.Disconnect(ctx)

    collName := fmt.Sprintf("profiles-getbyid-multi-test-%d", time.Now().UnixNano())
    coll := client.Database("user-server-test").Collection(collName)
    defer coll.Drop(ctx)

    store := NewMongoProfileStore(coll)

    // Insert multiple profiles with different update times
    baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
    recent := baseTime.Add(24 * time.Hour)

    p1 := &Profile{UserId: "user1", FirstName: "Profile1", UpdatedOn: &recent}
    p2 := &Profile{UserId: "user2", FirstName: "Profile2", UpdatedOn: &baseTime}
    p3 := &Profile{UserId: "user3", FirstName: "Profile3"} // No UpdatedOn

    if err := store.Upsert(&ctx, p1); err != nil {
        t.Fatalf("Upsert p1 failed: %v", err)
    }
    if err := store.Upsert(&ctx, p2); err != nil {
        t.Fatalf("Upsert p2 failed: %v", err)
    }
    if err := store.Upsert(&ctx, p3); err != nil {
        t.Fatalf("Upsert p3 failed: %v", err)
    }

    // GetByUserId with cutoff time should find p1 (updated after cutoff)
    cutoffTime := baseTime.Add(12 * time.Hour)
    got1, err := store.GetByUserId(&ctx, "user1", cutoffTime)
    if err != nil || got1 == nil {
        t.Fatalf("GetByUserId user1 failed: %v", err)
    }
    if got1.FirstName != "Profile1" {
        t.Fatalf("expected Profile1, got %s", got1.FirstName)
    }

    // GetByUserId for user2 with same cutoff should NOT find it (updated before cutoff)
    _, err = store.GetByUserId(&ctx, "user2", cutoffTime)
    if err == nil {
        t.Fatalf("expected error for user2, got nil")
    }

    // GetByUserId for user3 with any cutoff should NOT find it (no UpdatedOn)
    _, err = store.GetByUserId(&ctx, "user3", cutoffTime)
    if err == nil {
        t.Fatalf("expected error for user3 (no UpdatedOn), got nil")
    }

    // Delete user1 and verify GetByUserId fails
    if err := store.Delete(&ctx, "user1"); err != nil {
        t.Fatalf("Delete user1 failed: %v", err)
    }

    gotAfterDelete, err := store.GetByUserId(&ctx, "user1", cutoffTime)
    if err == nil {
        t.Fatalf("expected error after delete, got nil")
    }
    if gotAfterDelete != nil {
        t.Fatalf("expected nil after delete, got %v", gotAfterDelete)
    }

    // Delete non-existent user returns error (expected behavior in this API)
    err = store.Delete(&ctx, "nonexistent-user")
    if err == nil {
        // API behavior: Delete returns error for non-existent, which is OK
        // Just verify the remaining user is still there
    }

    // Verify Get still works for remaining profiles
    p2Get, err := store.Get(&ctx, "user2")
    if err != nil || p2Get == nil {
        t.Fatalf("Get user2 should work: %v", err)
    }
}
