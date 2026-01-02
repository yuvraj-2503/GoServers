// +build integration

package db

import (
	"context"
	mongodb "mongo-utils"
	"testing"
	"time"
)

// Integration tests for MongoProfileStore with real MongoDB
// Run with: go test -tags=integration ./user-server/profile/db -v

func TestMongoProfileStoreIntegration_Upsert(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	config := &mongodb.MongoConfig{
		ConnectionString: "mongodb://localhost:27017",
		Database:         "user-server-test",
		Username:         "",
		Password:         "",
	}

	coll, err := config.GetCollection("profile-integration-test")
	if err != nil {
		t.Fatalf("failed to get collection: %v", err)
	}
	defer coll.Drop(ctx)

	store := NewMongoProfileStore(coll)

	now := time.Now()
	profile := &Profile{
		UserId:    "user-123",
		FirstName: "John",
		LastName:  "Doe",
		UpdatedOn: &now,
	}

	err = store.Upsert(&ctx, profile)
	if err != nil {
		t.Errorf("Upsert() error = %v, wantErr false", err)
	}
}

func TestMongoProfileStoreIntegration_Get(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	config := &mongodb.MongoConfig{
		ConnectionString: "mongodb://localhost:27017",
		Database:         "user-server-test",
		Username:         "",
		Password:         "",
	}

	coll, err := config.GetCollection("profile-integration-test")
	if err != nil {
		t.Fatalf("failed to get collection: %v", err)
	}
	defer coll.Drop(ctx)

	store := NewMongoProfileStore(coll)

	now := time.Now()
	profile := &Profile{
		UserId:    "user-123",
		FirstName: "John",
		LastName:  "Doe",
		UpdatedOn: &now,
	}

	err = store.Upsert(&ctx, profile)
	if err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}

	got, err := store.Get(&ctx, "user-123")
	if err != nil {
		t.Errorf("Get() error = %v, wantErr false", err)
		return
	}

	if got == nil {
		t.Error("Get() returned nil, expected Profile")
		return
	}

	if got.UserId != "user-123" || got.FirstName != "John" {
		t.Errorf("Get() got = %v, want UserId=user-123 FirstName=John", got)
	}
}

func TestMongoProfileStoreIntegration_GetByUserId(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	config := &mongodb.MongoConfig{
		ConnectionString: "mongodb://localhost:27017",
		Database:         "user-server-test",
		Username:         "",
		Password:         "",
	}

	coll, err := config.GetCollection("profile-integration-test")
	if err != nil {
		t.Fatalf("failed to get collection: %v", err)
	}
	defer coll.Drop(ctx)

	store := NewMongoProfileStore(coll)

	now := time.Now()
	profile := &Profile{
		UserId:    "user-123",
		FirstName: "John",
		LastName:  "Doe",
		UpdatedOn: &now,
	}

	err = store.Upsert(&ctx, profile)
	if err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}

	// Fetch with a time threshold before the update
	beforeTime := now.Add(-1 * time.Hour)
	got, err := store.GetByUserId(&ctx, "user-123", beforeTime)

	if err != nil {
		t.Errorf("GetByUserId() error = %v, wantErr false", err)
		return
	}

	if got == nil {
		t.Error("GetByUserId() returned nil, expected Profile")
		return
	}

	if got.UserId != "user-123" {
		t.Errorf("GetByUserId() got = %v, want UserId=user-123", got)
	}
}

func TestMongoProfileStoreIntegration_Delete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	config := &mongodb.MongoConfig{
		ConnectionString: "mongodb://localhost:27017",
		Database:         "user-server-test",
		Username:         "",
		Password:         "",
	}

	coll, err := config.GetCollection("profile-integration-test")
	if err != nil {
		t.Fatalf("failed to get collection: %v", err)
	}
	defer coll.Drop(ctx)

	store := NewMongoProfileStore(coll)

	now := time.Now()
	profile := &Profile{
		UserId:    "user-123",
		FirstName: "John",
		LastName:  "Doe",
		UpdatedOn: &now,
	}

	err = store.Upsert(&ctx, profile)
	if err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}

	err = store.Delete(&ctx, "user-123")
	if err != nil {
		t.Errorf("Delete() error = %v, wantErr false", err)
		return
	}

	got, err := store.Get(&ctx, "user-123")
	if got != nil {
		t.Error("Delete() failed - profile still exists")
	}
}
