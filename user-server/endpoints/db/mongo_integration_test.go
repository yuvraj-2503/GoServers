// +build integration

package db

import (
	"context"
	mongodb "mongo-utils"
	"testing"
	"time"
)

// Integration tests for UrlMongoStore with real MongoDB
// Run with: go test -tags=integration ./user-server/endpoints/db -v

func TestUrlMongoStoreIntegration_Upsert(t *testing.T) {
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

	coll, err := config.GetCollection("urls_integration_test")
	if err != nil {
		t.Fatalf("failed to get collection: %v", err)
	}
	defer coll.Drop(ctx)

	store := NewUrlMongoStore(coll)

	now := time.Now()
	urls := []*UrlData{
		{
			Key:       "test-service",
			Url:       "http://localhost:3000",
			Env:       "TEST",
			UpdatedAt: &now,
		},
	}

	err = store.Upsert(&ctx, urls)
	if err != nil {
		t.Errorf("Upsert() error = %v, wantErr false", err)
	}
}

func TestUrlMongoStoreIntegration_Get(t *testing.T) {
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

	coll, err := config.GetCollection("urls_integration_test")
	if err != nil {
		t.Fatalf("failed to get collection: %v", err)
	}
	defer coll.Drop(ctx)

	store := NewUrlMongoStore(coll)

	now := time.Now()
	urls := []*UrlData{
		{
			Key:       "test-service",
			Url:       "http://localhost:3000",
			Env:       "TEST",
			UpdatedAt: &now,
		},
	}

	err = store.Upsert(&ctx, urls)
	if err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}

	got, err := store.Get(&ctx, "test-service", "TEST")
	if err != nil {
		t.Errorf("Get() error = %v, wantErr false", err)
		return
	}

	if got == nil {
		t.Error("Get() returned nil, expected UrlData")
		return
	}

	if got.Key != "test-service" || got.Env != "TEST" {
		t.Errorf("Get() got = %v, want Key=test-service Env=TEST", got)
	}
}

func TestUrlMongoStoreIntegration_GetAll(t *testing.T) {
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

	coll, err := config.GetCollection("urls_integration_test")
	if err != nil {
		t.Fatalf("failed to get collection: %v", err)
	}
	defer coll.Drop(ctx)

	store := NewUrlMongoStore(coll)

	now := time.Now()
	urls := []*UrlData{
		{
			Key:       "service1",
			Url:       "http://localhost:3000",
			Env:       "TEST",
			UpdatedAt: &now,
		},
		{
			Key:       "service2",
			Url:       "http://localhost:4000",
			Env:       "TEST",
			UpdatedAt: &now,
		},
	}

	err = store.Upsert(&ctx, urls)
	if err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}

	got, err := store.GetAll(&ctx, "TEST")
	if err != nil {
		t.Errorf("GetAll() error = %v, wantErr false", err)
		return
	}

	if len(got) != 2 {
		t.Errorf("GetAll() returned %d items, expected 2", len(got))
	}
}

func TestUrlMongoStoreIntegration_Delete(t *testing.T) {
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

	coll, err := config.GetCollection("urls_integration_test")
	if err != nil {
		t.Fatalf("failed to get collection: %v", err)
	}
	defer coll.Drop(ctx)

	store := NewUrlMongoStore(coll)

	now := time.Now()
	urls := []*UrlData{
		{
			Key:       "test-service",
			Url:       "http://localhost:3000",
			Env:       "TEST",
			UpdatedAt: &now,
		},
	}

	err = store.Upsert(&ctx, urls)
	if err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}

	err = store.Delete(&ctx, "test-service", "TEST")
	if err != nil {
		t.Errorf("Delete() error = %v, wantErr false", err)
	}

	got, err := store.Get(&ctx, "test-service", "TEST")
	if err == nil && got != nil {
		t.Error("Delete() failed - item still exists")
	}
}
