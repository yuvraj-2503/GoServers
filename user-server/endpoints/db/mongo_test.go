package db

import (
	"context"
	"fmt"
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
	type fields struct {
		urlColl *mongo.Collection
	}
	type args struct {
		ctx *context.Context
		key string
		env string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UrlMongoStore{
				urlColl: tt.fields.urlColl,
			}
			if err := u.Delete(tt.args.ctx, tt.args.key, tt.args.env); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUrlMongoStore_Get(t *testing.T) {
	type fields struct {
		urlColl *mongo.Collection
	}
	type args struct {
		ctx *context.Context
		key string
		env string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *UrlData
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UrlMongoStore{
				urlColl: tt.fields.urlColl,
			}
			got, err := u.Get(tt.args.ctx, tt.args.key, tt.args.env)
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
	type fields struct {
		urlColl *mongo.Collection
	}
	type args struct {
		ctx *context.Context
		env string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*UrlData
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Get All",
			fields: fields{
				urlColl: urlColl,
			},
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
			u := &UrlMongoStore{
				urlColl: tt.fields.urlColl,
			}
			got, err := u.GetAll(tt.args.ctx, tt.args.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() got = %v, want %v", got, tt.want)
				fmt.Println(got)
				for _, url := range got {
					fmt.Println(url)
				}
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
	type fields struct {
		urlColl *mongo.Collection
	}
	type args struct {
		ctx  *context.Context
		urls []*UrlData
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
				urlColl: urlColl,
			},
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
			u := &UrlMongoStore{
				urlColl: tt.fields.urlColl,
			}
			if err := u.Upsert(tt.args.ctx, tt.args.urls); (err != nil) != tt.wantErr {
				t.Errorf("Upsert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
