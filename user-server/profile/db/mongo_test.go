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

var profileColl *mongo.Collection
var ctx context.Context
var currTime = time.Date(2024, time.Month(11), 28, 0, 30, 30, 0, time.UTC)

func init() {
	mongoConfig := mongodb.MongoConfig{
		ConnectionString: "mongodb://localhost:27017",
		Database:         "user-server",
		Username:         "",
		Password:         "",
	}
	profileColl, _ = mongoConfig.GetCollection("profile-collection")
	ctx = context.Background()
}

// MockProfileStore is a mock implementation of ProfileStore for testing
type MockProfileStore struct {
	GetFunc       func(ctx *context.Context, userId string) (*Profile, error)
	GetByUserFunc func(ctx *context.Context, userId string, time time.Time) (*Profile, error)
	UpsertFunc    func(ctx *context.Context, profile *Profile) error
	DeleteFunc    func(ctx *context.Context, userId string) error
}

func (m *MockProfileStore) Get(ctx *context.Context, userId string) (*Profile, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, userId)
	}
	return nil, nil
}

func (m *MockProfileStore) GetByUserId(ctx *context.Context, userId string, t time.Time) (*Profile, error) {
	if m.GetByUserFunc != nil {
		return m.GetByUserFunc(ctx, userId, t)
	}
	return nil, nil
}

func (m *MockProfileStore) Upsert(ctx *context.Context, profile *Profile) error {
	if m.UpsertFunc != nil {
		return m.UpsertFunc(ctx, profile)
	}
	return nil
}

func (m *MockProfileStore) Delete(ctx *context.Context, userId string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, userId)
	}
	return nil
}

func TestMongoProfileStore_Delete(t *testing.T) {
	type args struct {
		ctx    *context.Context
		userId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Delete",
			args: args{
				ctx:    &ctx,
				userId: "6719252c58da11805939fea3",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockProfileStore{
				DeleteFunc: func(ctx *context.Context, userId string) error {
					return nil
				},
			}
			if err := mockStore.Delete(tt.args.ctx, tt.args.userId); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoProfileStore_Get(t *testing.T) {
	type args struct {
		ctx    *context.Context
		userId string
	}
	tests := []struct {
		name    string
		args    args
		want    *Profile
		wantErr bool
	}{
		{
			name: "Get",
			args: args{
				ctx:    &ctx,
				userId: "6719252c58da11805939fea3",
			},
			want: &Profile{
				UserId:    "6719252c58da11805939fea3",
				FirstName: "Yuvraj",
				LastName:  "Singh Rajpoot",
				UpdatedOn: &currTime,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockProfileStore{
				GetFunc: func(ctx *context.Context, userId string) (*Profile, error) {
					if userId == "6719252c58da11805939fea3" {
						return &Profile{
							UserId:    "6719252c58da11805939fea3",
							FirstName: "Yuvraj",
							LastName:  "Singh Rajpoot",
							UpdatedOn: &currTime,
						}, nil
					}
					return nil, nil
				},
			}
			got, err := mockStore.Get(tt.args.ctx, tt.args.userId)
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

func TestMongoProfileStore_Upsert(t *testing.T) {
	type args struct {
		ctx     *context.Context
		profile *Profile
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
				profile: &Profile{
					UserId:    "6719252c58da11805939fea3",
					FirstName: "Yuvraj",
					LastName:  "Singh Rajpoot",
					UpdatedOn: &currTime,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockProfileStore{
				UpsertFunc: func(ctx *context.Context, profile *Profile) error {
					return nil
				},
			}
			if err := mockStore.Upsert(tt.args.ctx, tt.args.profile); (err != nil) != tt.wantErr {
				t.Errorf("Upsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewMongoProfileStore(t *testing.T) {
	type args struct {
		profileColl *mongo.Collection
	}
	tests := []struct {
		name string
		args args
		want *MongoProfileStore
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMongoProfileStore(tt.args.profileColl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMongoProfileStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getUpdates(t *testing.T) {
	type args struct {
		profile *Profile
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
			if got := getUpdates(tt.args.profile); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getUpdates() = %v, want %v", got, tt.want)
			}
		})
	}
}
