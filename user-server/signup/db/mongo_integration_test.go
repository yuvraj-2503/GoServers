// +build integration

package db

import (
"context"
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/mongo/options"
"testing"
"user-server/common"
)

// Integration tests for MongoUserStore with real MongoDB
// Run with: go test -tags=integration ./user-server/signup/db -v

func getMongoClient(ctx context.Context) (*mongo.Client, error) {
	clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
	return mongo.Connect(ctx, clientOpts)
}

func TestMongoUserStoreIntegration_Insert(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	client, err := getMongoClient(ctx)
	if err != nil {
		t.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	coll := client.Database("user-server-test").Collection("users-integration-test")
	defer coll.Drop(ctx)

	store := NewMongoUserStore(coll)

	user := &User{
		UserId:  "user-123",
		EmailId: "test@example.com",
		PhoneNumber: &common.PhoneNumber{
			CountryCode: "+1",
			Number:      "1234567890",
		},
	}

	err = store.Insert(&ctx, user)
	if err != nil {
		t.Errorf("Insert() error = %v, wantErr false", err)
	}
}

func TestMongoUserStoreIntegration_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	client, err := getMongoClient(ctx)
	if err != nil {
		t.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	coll := client.Database("user-server-test").Collection("users-integration-test")
	defer coll.Drop(ctx)

	store := NewMongoUserStore(coll)

	user := &User{
		UserId:  "user-123",
		EmailId: "test@example.com",
	}

	err = store.Insert(&ctx, user)
	if err != nil {
		t.Fatalf("Insert() error = %v", err)
	}

	got, err := store.Get(&ctx, Filter{Key: UserId, Value: "user-123"})
	if err != nil {
		t.Errorf("Get() error = %v, wantErr false", err)
		return
	}

	if got == nil {
		t.Error("Get() returned nil, expected User")
		return
	}

	if got.UserId != "user-123" {
		t.Errorf("Get() got = %v, want UserId=user-123", got)
	}
}

func TestMongoUserStoreIntegration_GetByPhoneNumber(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	client, err := getMongoClient(ctx)
	if err != nil {
		t.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	coll := client.Database("user-server-test").Collection("users-integration-test")
	defer coll.Drop(ctx)

	store := NewMongoUserStore(coll)

	phoneNumber := &common.PhoneNumber{
		CountryCode: "+1",
		Number:      "1234567890",
	}

	user := &User{
		UserId:      "user-123",
		EmailId:     "test@example.com",
		PhoneNumber: phoneNumber,
	}

	err = store.Insert(&ctx, user)
	if err != nil {
		t.Fatalf("Insert() error = %v", err)
	}

	got, err := store.GetByPhoneNumber(&ctx, phoneNumber)
	if err != nil {
		t.Errorf("GetByPhoneNumber() error = %v, wantErr false", err)
		return
	}

	if got == nil {
		t.Error("GetByPhoneNumber() returned nil, expected User")
		return
	}

	if got.UserId != "user-123" {
		t.Errorf("GetByPhoneNumber() got = %v, want UserId=user-123", got)
	}
}

func TestMongoUserStoreIntegration_CheckExists(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	client, err := getMongoClient(ctx)
	if err != nil {
		t.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	coll := client.Database("user-server-test").Collection("users-integration-test")
	defer coll.Drop(ctx)

	store := NewMongoUserStore(coll)

	user := &User{
		UserId:  "user-123",
		EmailId: "test@example.com",
	}

	err = store.Insert(&ctx, user)
	if err != nil {
		t.Fatalf("Insert() error = %v", err)
	}

	exists, err := store.CheckExists(&ctx, Filter{Key: UserId, Value: "user-123"})
	if err != nil {
		t.Errorf("CheckExists() error = %v, wantErr false", err)
		return
	}

	if !exists {
		t.Error("CheckExists() returned false, expected true")
	}
}

func TestMongoUserStoreIntegration_UpdateEmailId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	client, err := getMongoClient(ctx)
	if err != nil {
		t.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	coll := client.Database("user-server-test").Collection("users-integration-test")
	defer coll.Drop(ctx)

	store := NewMongoUserStore(coll)

	user := &User{
		UserId:  "user-123",
		EmailId: "old@example.com",
	}

	err = store.Insert(&ctx, user)
	if err != nil {
		t.Fatalf("Insert() error = %v", err)
	}

	err = store.UpdateEmailId(&ctx, "user-123", "new@example.com")
	if err != nil {
		t.Errorf("UpdateEmailId() error = %v, wantErr false", err)
		return
	}

	got, err := store.Get(&ctx, Filter{Key: UserId, Value: "user-123"})
	if err != nil || got == nil {
		t.Errorf("Get() error = %v", err)
		return
	}

	if got.EmailId != "new@example.com" {
		t.Errorf("UpdateEmailId() failed - EmailId = %v, want new@example.com", got.EmailId)
	}
}

func TestMongoUserStoreIntegration_DeleteByUserId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	client, err := getMongoClient(ctx)
	if err != nil {
		t.Fatalf("failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(ctx)

	coll := client.Database("user-server-test").Collection("users-integration-test")
	defer coll.Drop(ctx)

	store := NewMongoUserStore(coll)

	user := &User{
		UserId:  "user-123",
		EmailId: "test@example.com",
	}

	err = store.Insert(&ctx, user)
	if err != nil {
		t.Fatalf("Insert() error = %v", err)
	}

	err = store.DeleteByUserId(&ctx, "user-123")
	if err != nil {
		t.Errorf("DeleteByUserId() error = %v, wantErr false", err)
		return
	}

	exists, err := store.CheckExists(&ctx, Filter{Key: UserId, Value: "user-123"})
	if exists {
		t.Error("DeleteByUserId() failed - user still exists")
	}
}
