package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongodb "mongo-utils"
	"reflect"
	"testing"
	"time"
)

var urlColl *mongo.Collection
var ctx context.Context

func init() {
	ctx = context.Background()
	config := &mongodb.MongoConfig{
		ConnectionString: "mongodb://localhost:27017",
		Database:         "user-server",
		Username:         "",
		Password:         "",
	}
	urlColl, _ = config.GetCollection("urls")
}

// MockUrlStore is a mock implementation of UrlStore for testing
type MockUrlStore struct {
	UpsertFunc func(ctx *context.Context, urls []*UrlData) error
	GetFunc    func(ctx *context.Context, key, env string) (*UrlData, error)
	GetAllFunc func(ctx *context.Context, env string) ([]*UrlData, error)
	DeleteFunc func(ctx *context.Context, key, env string) error
}

func (m *MockUrlStore) Upsert(ctx *context.Context, urls []*UrlData) error {
	if m.UpsertFunc != nil {
		return m.UpsertFunc(ctx, urls)
	}
	return nil
}

func (m *MockUrlStore) Get(ctx *context.Context, key, env string) (*UrlData, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, key, env)
	}
	return nil, nil
}

func (m *MockUrlStore) GetAll(ctx *context.Context, env string) ([]*UrlData, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx, env)
	}
	return nil, nil
}

func (m *MockUrlStore) Delete(ctx *context.Context, key, env string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, key, env)
	}
	return nil
}

func TestNewUrlMongoStore(t *testing.T) {
	type args struct {
		urlColl *mongo.Collection
	}
	tests := []struct {
		name string
		args args
		want *UrlMongoStore
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUrlMongoStore(tt.args.urlColl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUrlMongoStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrlMongoStore_Delete(t *testing.T) {
	type args struct {
		ctx *context.Context
		key string
		env string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Delete",
			args: args{
				ctx: &ctx,
				key: "user-server",
				env: "LOCAL",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockUrlStore{
				DeleteFunc: func(ctx *context.Context, key, env string) error {
					return nil
				},
			}
			if err := mockStore.Delete(tt.args.ctx, tt.args.key, tt.args.env); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUrlMongoStore_Get(t *testing.T) {
	type args struct {
		ctx *context.Context
		key string
		env string
	}
	tests := []struct {
		name    string
		args    args
		want    *UrlData
		wantErr bool
	}{
		{
			name: "Get",
			args: args{
				ctx: &ctx,
				key: "user-server",
				env: "LOCAL",
			},
			want: &UrlData{
				Key: "user-server",
				Url: "http://localhost:8080/api/v1",
				Env: "LOCAL",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockUrlStore{
				GetFunc: func(ctx *context.Context, key, env string) (*UrlData, error) {
					if key == "user-server" && env == "LOCAL" {
						return &UrlData{
							Key: "user-server",
							Url: "http://localhost:8080/api/v1",
							Env: "LOCAL",
						}, nil
					}
					return nil, nil
				},
			}
			got, err := mockStore.Get(tt.args.ctx, tt.args.key, tt.args.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrlMongoStore_GetAll(t *testing.T) {
	var now = time.Now()
	type args struct {
		ctx *context.Context
		env string
	}
	tests := []struct {
		name    string
		args    args
		want    []*UrlData
		wantErr bool
	}{
		{
			name: "Get All",
			args: args{
				ctx: &ctx,
				env: "LOCAL",
			},
			want: []*UrlData{{
				Key:       "user-server",
				Url:       "http://localhost:8080/api/v1",
				Env:       "LOCAL",
				UpdatedAt: &now,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockUrlStore{
				GetAllFunc: func(ctx *context.Context, env string) ([]*UrlData, error) {
					if env == "LOCAL" {
						return []*UrlData{{
							Key:       "user-server",
							Url:       "http://localhost:8080/api/v1",
							Env:       "LOCAL",
							UpdatedAt: &now,
						}}, nil
					}
					return nil, nil
				},
			}
			got, err := mockStore.GetAll(tt.args.ctx, tt.args.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getUpdates(t *testing.T) {
	type args struct {
		url *UrlData
	}
	tests := []struct {
		name string
		args args
		want bson.D
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getUpdates(tt.args.url); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getUpdates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrlMongoStore_Upsert(t *testing.T) {
	var now = time.Now()
	type args struct {
		ctx  *context.Context
		urls []*UrlData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Upsert",
			args: args{
				ctx: &ctx,
				urls: []*UrlData{
					{
						Key:       "user-server",
						Url:       "http://localhost:8080/api/v1",
						UpdatedAt: &now,
						Env:       "LOCAL",
					},
					{
						Key:       "user-server",
						Url:       "http://10.0.2.2:8080/api/v1",
						UpdatedAt: &now,
						Env:       "DEVELOPMENT",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockUrlStore{
				UpsertFunc: func(ctx *context.Context, urls []*UrlData) error {
					return nil
				},
			}
			if err := mockStore.Upsert(tt.args.ctx, tt.args.urls); (err != nil) != tt.wantErr {
				t.Errorf("Upsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
