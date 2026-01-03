package db

import (
    "context"
    "fmt"
    "testing"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func getMongoClientLocal(ctx context.Context) (*mongo.Client, error) {
    clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
    return mongo.Connect(ctx, clientOpts)
}

func TestUrlMongoStore_Unit_UpsertGetGetAllDelete(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration-style unit test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("failed to connect to MongoDB: %v", err)
    }
    defer client.Disconnect(ctx)

    collName := fmt.Sprintf("urls-unit-test-%d", time.Now().UnixNano())
    coll := client.Database("user-server-test").Collection(collName)
    defer coll.Drop(ctx)

    store := NewUrlMongoStore(coll)

    now := time.Now()
    urls := []*UrlData{
        {Key: "k1", Url: "https://a.example", Env: "dev", UpdatedAt: &now},
        {Key: "k2", Url: "https://b.example", Env: "dev", UpdatedAt: &now},
    }

    if err := store.Upsert(&ctx, urls); err != nil {
        t.Fatalf("Upsert() error = %v", err)
    }

    // Get k1
    got, err := store.Get(&ctx, "k1", "dev")
    if err != nil {
        t.Fatalf("Get() error = %v", err)
    }
    if got == nil || got.Url != "https://a.example" {
        t.Fatalf("unexpected Get result: %+v", got)
    }

    // GetAll
    all, err := store.GetAll(&ctx, "dev")
    if err != nil {
        t.Fatalf("GetAll() error = %v", err)
    }
    if len(all) < 2 {
        t.Fatalf("expected at least 2 entries, got %d", len(all))
    }

    // Delete one
    if err := store.Delete(&ctx, "k1", "dev"); err != nil {
        t.Fatalf("Delete() error = %v", err)
    }

    // Get should now return not found
    _, err = store.Get(&ctx, "k1", "dev")
    if err == nil {
        t.Fatalf("expected error after delete, got nil")
    }
}

func TestUrlMongoStore_Unit_GetAllEmpty_And_UpsertPartialFields(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration-style unit test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("failed to connect to MongoDB: %v", err)
    }
    defer client.Disconnect(ctx)

    collName := fmt.Sprintf("urls-empty-test-%d", time.Now().UnixNano())
    coll := client.Database("user-server-test").Collection(collName)
    defer coll.Drop(ctx)

    store := NewUrlMongoStore(coll)

    // GetAll on empty collection
    all, err := store.GetAll(&ctx, "empty-env")
    if err != nil {
        t.Logf("GetAll on empty returned error: %v (may be acceptable)", err)
    } else if len(all) > 0 {
        t.Fatalf("expected empty slice, got %d entries", len(all))
    }

    // Upsert with partial fields (only key and url)
    urls := []*UrlData{
        {Key: "k1", Url: "https://minimal.example", Env: "test"},
    }
    if err := store.Upsert(&ctx, urls); err != nil {
        t.Fatalf("Upsert with partial fields failed: %v", err)
    }

    // Verify the entry exists
    got, err := store.Get(&ctx, "k1", "test")
    if err != nil {
        t.Fatalf("Get after upsert failed: %v", err)
    }
    if got == nil || got.Url != "https://minimal.example" {
        t.Fatalf("unexpected Get result: %+v", got)
    }

    // Upsert same key again with updated URL
    urls[0].Url = "https://updated.example"
    if err := store.Upsert(&ctx, urls); err != nil {
        t.Fatalf("Upsert update failed: %v", err)
    }

    got2, err := store.Get(&ctx, "k1", "test")
    if err != nil || got2 == nil {
        t.Fatalf("Get after update failed: %v", err)
    }
    if got2.Url != "https://updated.example" {
        t.Fatalf("expected updated URL, got %s", got2.Url)
    }
}

func TestUrlMongoStore_Unit_DeleteNotFound_And_GetAllNoDocuments(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration-style unit test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("failed to connect to MongoDB: %v", err)
    }
    defer client.Disconnect(ctx)

    collName := fmt.Sprintf("urls-delete-test-%d", time.Now().UnixNano())
    coll := client.Database("user-server-test").Collection(collName)
    defer coll.Drop(ctx)

    store := NewUrlMongoStore(coll)

    // Delete non-existent (empty collection) - should succeed with no error (doesn't check DeletedCount)
    if err := store.Delete(&ctx, "k-missing", "env-missing"); err != nil {
        t.Logf("Delete on non-existent returned error (may be acceptable): %v", err)
    }

    // Insert one url
    urls := []*UrlData{{Key: "k1", Url: "https://test.example", Env: "test"}}
    if err := store.Upsert(&ctx, urls); err != nil {
        t.Fatalf("Upsert failed: %v", err)
    }

    // Delete it
    if err := store.Delete(&ctx, "k1", "test"); err != nil {
        t.Fatalf("Delete existing failed: %v", err)
    }

    // Delete again (should also succeed per implementation)
    if err := store.Delete(&ctx, "k1", "test"); err != nil {
        t.Logf("Delete on already-deleted returned error (may be acceptable): %v", err)
    }

    // GetAll after delete
    all, err := store.GetAll(&ctx, "test")
    if err != nil {
        t.Logf("GetAll after delete returned error: %v (may be acceptable)", err)
    } else if len(all) > 0 {
        t.Fatalf("expected empty GetAll result, got %d entries", len(all))
    }
}

func TestUrlMongoStore_Unit_GetNotFound_And_GetAllMultipleEnvs(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration-style unit test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("failed to connect to MongoDB: %v", err)
    }
    defer client.Disconnect(ctx)

    collName := fmt.Sprintf("urls-multienv-test-%d", time.Now().UnixNano())
    coll := client.Database("user-server-test").Collection(collName)
    defer coll.Drop(ctx)

    store := NewUrlMongoStore(coll)

    // Get from empty collection
    _, err = store.Get(&ctx, "k-missing", "env-missing")
    if err == nil {
        t.Fatalf("expected NotFoundError on get from empty, got nil")
    }

    // Upsert to multiple environments
    urls := []*UrlData{
        {Key: "shared-key", Url: "https://dev.example", Env: "dev"},
        {Key: "shared-key", Url: "https://prod.example", Env: "prod"},
        {Key: "other-key", Url: "https://test.example", Env: "dev"},
    }
    if err := store.Upsert(&ctx, urls); err != nil {
        t.Fatalf("Upsert failed: %v", err)
    }

    // GetAll for dev env should return 2 entries
    devAll, err := store.GetAll(&ctx, "dev")
    if err != nil {
        t.Fatalf("GetAll(dev) error: %v", err)
    }
    if len(devAll) != 2 {
        t.Fatalf("expected 2 dev entries, got %d", len(devAll))
    }

    // GetAll for prod env should return 1 entry
    prodAll, err := store.GetAll(&ctx, "prod")
    if err != nil {
        t.Fatalf("GetAll(prod) error: %v", err)
    }
    if len(prodAll) != 1 {
        t.Fatalf("expected 1 prod entry, got %d", len(prodAll))
    }

    // Get specific key/env combo
    got, err := store.Get(&ctx, "shared-key", "prod")
    if err != nil || got == nil {
        t.Fatalf("Get shared-key/prod failed: %v", err)
    }
    if got.Url != "https://prod.example" {
        t.Fatalf("expected prod URL, got %s", got.Url)
    }

    // Get from non-matching env
    _, err = store.Get(&ctx, "shared-key", "staging")
    if err == nil {
        t.Fatalf("expected NotFoundError for non-matching env, got nil")
    }
}

func TestUrlMongoStore_Unit_Delete_With_MultipleFilters_And_Empty_Collection(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping Mongo-backed test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("connect error: %v", err)
    }
    defer client.Disconnect(ctx)

    collName := fmt.Sprintf("urls-delete-multi-test-%d", time.Now().UnixNano())
    coll := client.Database("user-server-test").Collection(collName)
    defer coll.Drop(ctx)

    store := NewUrlMongoStore(coll)

    // Test 1: Delete from empty collection (should pass, no error even if nothing deleted)
    err = store.Delete(&ctx, "nonexistent-key", "prod")
    if err != nil {
        t.Fatalf("Delete from empty collection should not error: %v", err)
    }

    // Insert some URLs
    now := time.Now()
    u1 := &UrlData{Key: "key1", Env: "prod", Url: "https://api1.example.com", UpdatedAt: &now}
    u2 := &UrlData{Key: "key1", Env: "dev", Url: "https://api1-dev.example.com", UpdatedAt: &now}
    u3 := &UrlData{Key: "key2", Env: "prod", Url: "https://api2.example.com", UpdatedAt: &now}

    if err := store.Upsert(&ctx, []*UrlData{u1}); err != nil {
        t.Fatalf("Upsert u1 failed: %v", err)
    }
    if err := store.Upsert(&ctx, []*UrlData{u2}); err != nil {
        t.Fatalf("Upsert u2 failed: %v", err)
    }
    if err := store.Upsert(&ctx, []*UrlData{u3}); err != nil {
        t.Fatalf("Upsert u3 failed: %v", err)
    }

    // Delete key1 from prod (should remove u1, keep u2 and u3)
    if err := store.Delete(&ctx, "key1", "prod"); err != nil {
        t.Fatalf("Delete by key/env failed: %v", err)
    }

    // Verify u2 still exists (key1/dev)
    got, err := store.Get(&ctx, "key1", "dev")
    if err != nil || got == nil {
        t.Fatalf("expected to find key1/dev, got: %v", err)
    }

    // Verify key2/prod still exists
    got2, err := store.Get(&ctx, "key2", "prod")
    if err != nil || got2 == nil {
        t.Fatalf("expected to find key2/prod, got: %v", err)
    }

    // Test 2: Delete non-existent (key/env combination)
    err = store.Delete(&ctx, "key3", "staging")
    if err != nil {
        t.Fatalf("Delete non-existent should not error: %v", err)
    }

    // Verify remaining URLs unchanged
    all, err := store.GetAll(&ctx, "prod")
    if err != nil {
        t.Fatalf("GetAll error: %v", err)
    }
    if len(all) != 1 {
        t.Fatalf("expected 1 prod URL remaining, got %d", len(all))
    }

    allDev, err := store.GetAll(&ctx, "dev")
    if err != nil {
        t.Fatalf("GetAll dev error: %v", err)
    }
    if len(allDev) != 1 {
        t.Fatalf("expected 1 dev URL remaining, got %d", len(allDev))
    }
}

func TestUrlMongoStore_Unit_Upsert_Partial_Fields_And_GetAll_Variants(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping Mongo-backed test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("connect error: %v", err)
    }
    defer client.Disconnect(ctx)

    collName := fmt.Sprintf("urls-upsert-partial-test-%d", time.Now().UnixNano())
    coll := client.Database("user-server-test").Collection(collName)
    defer coll.Drop(ctx)

    store := NewUrlMongoStore(coll)

    now := time.Now()

    // Insert initial URL
    u1 := &UrlData{Key: "k1", Env: "prod", Url: "https://api1.example.com", UpdatedAt: &now}
    if err := store.Upsert(&ctx, []*UrlData{u1}); err != nil {
        t.Fatalf("Initial upsert failed: %v", err)
    }

    // Upsert same key with different URL (should update)
    u1Updated := &UrlData{Key: "k1", Env: "prod", Url: "https://api1-new.example.com", UpdatedAt: &now}
    if err := store.Upsert(&ctx, []*UrlData{u1Updated}); err != nil {
        t.Fatalf("Upsert update failed: %v", err)
    }

    // Verify update
    got, err := store.Get(&ctx, "k1", "prod")
    if err != nil || got == nil {
        t.Fatalf("Get after upsert failed: %v", err)
    }
    if got.Url != "https://api1-new.example.com" {
        t.Fatalf("expected updated URL, got %s", got.Url)
    }

    // GetAll with multiple records
    u2 := &UrlData{Key: "k2", Env: "prod", Url: "https://api2.example.com", UpdatedAt: &now}
    u3 := &UrlData{Key: "k3", Env: "staging", Url: "https://api3-staging.example.com", UpdatedAt: &now}
    _ = store.Upsert(&ctx, []*UrlData{u2})
    _ = store.Upsert(&ctx, []*UrlData{u3})

    // GetAll for prod should return 2
    all, err := store.GetAll(&ctx, "prod")
    if err != nil {
        t.Fatalf("GetAll prod failed: %v", err)
    }
    if len(all) != 2 {
        t.Fatalf("expected 2 prod URLs, got %d", len(all))
    }

    // GetAll for staging should return 1
    allStaging, err := store.GetAll(&ctx, "staging")
    if err != nil {
        t.Fatalf("GetAll staging failed: %v", err)
    }
    if len(allStaging) != 1 {
        t.Fatalf("expected 1 staging URL, got %d", len(allStaging))
    }

    // GetAll for non-existent env should return empty
    allNone, err := store.GetAll(&ctx, "nonexistent")
    if err != nil {
        t.Fatalf("GetAll nonexistent should not error: %v", err)
    }
    if len(allNone) != 0 {
        t.Fatalf("expected 0 URLs for nonexistent env, got %d", len(allNone))
    }
}

func TestUrlMongoStore_Unit_Get_Error_Scenarios_And_Delete_By_Key_Env(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping Mongo-backed test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("connect error: %v", err)
    }
    defer client.Disconnect(ctx)

    collName := fmt.Sprintf("urls-get-delete-test-%d", time.Now().UnixNano())
    coll := client.Database("user-server-test").Collection(collName)
    defer coll.Drop(ctx)

    store := NewUrlMongoStore(coll)

    // Get from empty collection - should error
    _, err = store.Get(&ctx, "nonexistent", "prod")
    if err == nil {
        t.Fatalf("Get from empty collection should error")
    }

    // Insert test data
    now := time.Now()
    u1 := &UrlData{Key: "key1", Env: "prod", Url: "https://prod.example.com", UpdatedAt: &now}
    u2 := &UrlData{Key: "key1", Env: "dev", Url: "https://dev.example.com", UpdatedAt: &now}
    _ = store.Upsert(&ctx, []*UrlData{u1})
    _ = store.Upsert(&ctx, []*UrlData{u2})

    // Get same key, different env
    got1, err := store.Get(&ctx, "key1", "prod")
    if err != nil || got1 == nil {
        t.Fatalf("Get key1 prod failed: %v", err)
    }
    if got1.Url != "https://prod.example.com" {
        t.Fatalf("expected prod URL, got %s", got1.Url)
    }

    got2, err := store.Get(&ctx, "key1", "dev")
    if err != nil || got2 == nil {
        t.Fatalf("Get key1 dev failed: %v", err)
    }
    if got2.Url != "https://dev.example.com" {
        t.Fatalf("expected dev URL, got %s", got2.Url)
    }

    // Get with wrong env - should error
    _, err = store.Get(&ctx, "key1", "staging")
    if err == nil {
        t.Fatalf("Get wrong env should error")
    }

    // Delete by key and env
    if err := store.Delete(&ctx, "key1", "prod"); err != nil {
        t.Fatalf("Delete key1 prod failed: %v", err)
    }

    // Verify prod deleted but dev still exists
    _, err = store.Get(&ctx, "key1", "prod")
    if err == nil {
        t.Fatalf("expected error after delete, got nil")
    }

    got3, err := store.Get(&ctx, "key1", "dev")
    if err != nil || got3 == nil {
        t.Fatalf("dev should still exist after deleting prod: %v", err)
    }

    // Delete non-existent - should not error
    if err := store.Delete(&ctx, "nonexistent", "prod"); err != nil {
        t.Fatalf("Delete non-existent should not error: %v", err)
    }
}
