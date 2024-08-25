package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"testing"
	"user-server/common"
)

func TestMongoUserStore_CheckExists(t *testing.T) {
	type fields struct {
		userColl *mongo.Collection
	}
	type args struct {
		ctx    *context.Context
		filter Filter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &MongoUserStore{
				userColl: tt.fields.userColl,
			}
			got, err := u.CheckExists(tt.args.ctx, tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckExists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoUserStore_CheckIfMobileExists(t *testing.T) {
	type fields struct {
		userColl *mongo.Collection
	}
	type args struct {
		ctx         *context.Context
		phoneNumber *common.PhoneNumber
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &MongoUserStore{
				userColl: tt.fields.userColl,
			}
			got, err := u.CheckIfMobileExists(tt.args.ctx, tt.args.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckIfMobileExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckIfMobileExists() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoUserStore_Delete(t *testing.T) {
	type fields struct {
		userColl *mongo.Collection
	}
	type args struct {
		ctx    *context.Context
		filter Filter
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
			u := &MongoUserStore{
				userColl: tt.fields.userColl,
			}
			if err := u.Delete(tt.args.ctx, tt.args.filter); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoUserStore_DeleteByUserId(t *testing.T) {
	type fields struct {
		userColl *mongo.Collection
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &MongoUserStore{
				userColl: tt.fields.userColl,
			}
			if err := u.DeleteByUserId(tt.args.ctx, tt.args.userId); (err != nil) != tt.wantErr {
				t.Errorf("DeleteByUserId() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoUserStore_Get(t *testing.T) {
	type fields struct {
		userColl *mongo.Collection
	}
	type args struct {
		ctx    *context.Context
		filter Filter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &MongoUserStore{
				userColl: tt.fields.userColl,
			}
			got, err := u.Get(tt.args.ctx, tt.args.filter)
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

func TestMongoUserStore_GetByPhoneNumber(t *testing.T) {
	type fields struct {
		userColl *mongo.Collection
	}
	type args struct {
		ctx         *context.Context
		phoneNumber *common.PhoneNumber
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &MongoUserStore{
				userColl: tt.fields.userColl,
			}
			got, err := u.GetByPhoneNumber(tt.args.ctx, tt.args.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetByPhoneNumber() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMongoUserStore_Insert(t *testing.T) {
	type fields struct {
		userColl *mongo.Collection
	}
	type args struct {
		ctx  *context.Context
		user *User
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
			u := &MongoUserStore{
				userColl: tt.fields.userColl,
			}
			if err := u.Insert(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoUserStore_UpdateEmailId(t *testing.T) {
	type fields struct {
		userColl *mongo.Collection
	}
	type args struct {
		ctx     *context.Context
		userId  string
		emailId string
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
			u := &MongoUserStore{
				userColl: tt.fields.userColl,
			}
			if err := u.UpdateEmailId(tt.args.ctx, tt.args.userId, tt.args.emailId); (err != nil) != tt.wantErr {
				t.Errorf("UpdateEmailId() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoUserStore_UpdatePhoneNumber(t *testing.T) {
	type fields struct {
		userColl *mongo.Collection
	}
	type args struct {
		ctx         *context.Context
		userId      string
		phoneNumber *common.PhoneNumber
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
			u := &MongoUserStore{
				userColl: tt.fields.userColl,
			}
			if err := u.UpdatePhoneNumber(tt.args.ctx, tt.args.userId, tt.args.phoneNumber); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewMongoUserStore(t *testing.T) {
	type args struct {
		collection *mongo.Collection
	}
	tests := []struct {
		name string
		args args
		want *MongoUserStore
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMongoUserStore(tt.args.collection); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMongoUserStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createBsonFilter(t *testing.T) {
	type args struct {
		filter Filter
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
			if got := createBsonFilter(tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createBsonFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
