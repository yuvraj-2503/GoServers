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

func TestMongoProfileStore_Delete(t *testing.T) {
	type fields struct {
		profileColl *mongo.Collection
	}
	type args struct {
		ctx    *context.Context
		userId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Delete",
			fields: fields{
				profileColl: profileColl,
			},
			args: args{
				ctx:    &ctx,
				userId: "6719252c58da11805939fea3",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &MongoProfileStore{
				profileColl: tt.fields.profileColl,
			}
			if err := p.Delete(tt.args.ctx, tt.args.userId); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoProfileStore_Get(t *testing.T) {
	type fields struct {
		profileColl *mongo.Collection
	}
	type args struct {
		ctx    *context.Context
		userId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Profile
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Get",
			fields: fields{
				profileColl: profileColl,
			},
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
			p := &MongoProfileStore{
				profileColl: tt.fields.profileColl,
			}
			got, err := p.Get(tt.args.ctx, tt.args.userId)
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
	//var now = time.Now()
	type fields struct {
		profileColl *mongo.Collection
	}
	type args struct {
		ctx     *context.Context
		profile *Profile
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Upsert",
			fields: fields{
				profileColl: profileColl,
			},
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
			p := &MongoProfileStore{
				profileColl: tt.fields.profileColl,
			}
			if err := p.Upsert(tt.args.ctx, tt.args.profile); (err != nil) != tt.wantErr {
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
