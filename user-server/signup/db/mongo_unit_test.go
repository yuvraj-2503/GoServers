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

func TestMongoUserStore_Unit_UpdateEmail_UpdatePhone_Delete_CheckIfMobileExists(t *testing.T) {
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

    phone := &common.PhoneNumber{CountryCode: "+91", Number: "9998887777"}
    user := &User{
        UserId:      "unit-user-2",
        EmailId:     "before@example.com",
        PhoneNumber: phone,
    }

    if err := store.Insert(&ctx, user); err != nil {
        t.Fatalf("Insert() error = %v", err)
    }

    // Update email
    if err := store.UpdateEmailId(&ctx, "unit-user-2", "after@example.com"); err != nil {
        t.Fatalf("UpdateEmailId() error = %v", err)
    }
    got, err := store.Get(&ctx, Filter{Key: UserId, Value: "unit-user-2"})
    if err != nil || got == nil {
        t.Fatalf("Get() error = %v, got=%v", err, got)
    }
    if got.EmailId != "after@example.com" {
        t.Fatalf("expected email updated, got %s", got.EmailId)
    }

    // Update phone
    newPhone := &common.PhoneNumber{CountryCode: "+44", Number: "7776665555"}
    if err := store.UpdatePhoneNumber(&ctx, "unit-user-2", newPhone); err != nil {
        t.Fatalf("UpdatePhoneNumber() error = %v", err)
    }
    got2, err := store.GetByPhoneNumber(&ctx, newPhone)
    if err != nil || got2 == nil {
        t.Fatalf("GetByPhoneNumber() error = %v, got=%v", err, got2)
    }
    if got2.UserId != "unit-user-2" {
        t.Fatalf("expected user id unit-user-2, got %s", got2.UserId)
    }

    // CheckIfMobileExists
    exists, err := store.CheckIfMobileExists(&ctx, newPhone)
    if err != nil {
        t.Fatalf("CheckIfMobileExists() error = %v", err)
    }
    if !exists {
        t.Fatalf("expected mobile to exist")
    }

    // Delete by user id
    if err := store.DeleteByUserId(&ctx, "unit-user-2"); err != nil {
        t.Fatalf("DeleteByUserId() error = %v", err)
    }
    existsAfter, _ := store.CheckExists(&ctx, Filter{Key: UserId, Value: "unit-user-2"})
    if existsAfter {
        t.Fatalf("expected user to be deleted")
    }
}

func TestMongoUserStore_Unit_InsertDuplicate_And_CheckExistsVariants(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration-style unit test in short mode")
    }

    ctx := context.Background()
    client, err := getMongoClientLocal(ctx)
    if err != nil {
        t.Fatalf("failed to connect: %v", err)
    }
    defer client.Disconnect(ctx)

    collName := fmt.Sprintf("users-unit-dupcheck-%d", time.Now().UnixNano())
    coll := client.Database("user-server-test").Collection(collName)
    defer coll.Drop(ctx)

    store := NewMongoUserStore(coll)

    // Insert a user
    u1 := &User{UserId: "dupcheck-user", EmailId: "test@example.com", PhoneNumber: &common.PhoneNumber{CountryCode: "+1", Number: "5551234567"}}
    if err := store.Insert(&ctx, u1); err != nil {
        t.Fatalf("Insert() error = %v", err)
    }

    // CheckExists by UserId
    exists, err := store.CheckExists(&ctx, Filter{Key: UserId, Value: "dupcheck-user"})
    if err != nil || !exists {
        t.Fatalf("CheckExists UserId failed: %v, %v", exists, err)
    }

    // CheckExists by EmailId
    exists, err = store.CheckExists(&ctx, Filter{Key: EmailId, Value: "test@example.com"})
    if err != nil || !exists {
        t.Fatalf("CheckExists EmailId failed: %v, %v", exists, err)
    }

    // CheckExists by non-existent EmailId
    exists, err = store.CheckExists(&ctx, Filter{Key: EmailId, Value: "notfound@example.com"})
    if err != nil || exists {
        t.Fatalf("CheckExists non-existent should be false: %v, %v", exists, err)
    }

    // CheckIfMobileExists
    phone := &common.PhoneNumber{CountryCode: "+1", Number: "5551234567"}
    exists, err = store.CheckIfMobileExists(&ctx, phone)
    if err != nil || !exists {
        t.Fatalf("CheckIfMobileExists existing phone failed: %v, %v", exists, err)
    }

    // CheckIfMobileExists non-existent
    nonExistPhone := &common.PhoneNumber{CountryCode: "+1", Number: "9999999999"}
    exists, err = store.CheckIfMobileExists(&ctx, nonExistPhone)
    if err != nil || exists {
        t.Fatalf("CheckIfMobileExists non-existent should be false: %v, %v", exists, err)
    }
}

