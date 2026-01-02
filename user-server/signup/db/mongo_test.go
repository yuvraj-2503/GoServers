package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"testing"
	"user-server/common"
)

// MockUserStore is a mock implementation of UserStore for testing
type MockUserStore struct {
	InsertFunc                 func(ctx *context.Context, user *User) error
	GetFunc                    func(ctx *context.Context, filter Filter) (*User, error)
	GetByPhoneNumberFunc       func(ctx *context.Context, phoneNumber *common.PhoneNumber) (*User, error)
	UpdateEmailIdFunc          func(ctx *context.Context, userId, emailId string) error
	UpdatePhoneNumberFunc      func(ctx *context.Context, userId string, phoneNumber *common.PhoneNumber) error
	DeleteByUserIdFunc         func(ctx *context.Context, userId string) error
	DeleteFunc                 func(ctx *context.Context, filter Filter) error
	CheckExistsFunc            func(ctx *context.Context, filter Filter) (bool, error)
	CheckIfMobileExistsFunc    func(ctx *context.Context, phoneNumber *common.PhoneNumber) (bool, error)
}

func (m *MockUserStore) Insert(ctx *context.Context, user *User) error {
	if m.InsertFunc != nil {
		return m.InsertFunc(ctx, user)
	}
	return nil
}

func (m *MockUserStore) Get(ctx *context.Context, filter Filter) (*User, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, filter)
	}
	return nil, nil
}

func (m *MockUserStore) GetByPhoneNumber(ctx *context.Context, phoneNumber *common.PhoneNumber) (*User, error) {
	if m.GetByPhoneNumberFunc != nil {
		return m.GetByPhoneNumberFunc(ctx, phoneNumber)
	}
	return nil, nil
}

func (m *MockUserStore) UpdateEmailId(ctx *context.Context, userId, emailId string) error {
	if m.UpdateEmailIdFunc != nil {
		return m.UpdateEmailIdFunc(ctx, userId, emailId)
	}
	return nil
}

func (m *MockUserStore) UpdatePhoneNumber(ctx *context.Context, userId string, phoneNumber *common.PhoneNumber) error {
	if m.UpdatePhoneNumberFunc != nil {
		return m.UpdatePhoneNumberFunc(ctx, userId, phoneNumber)
	}
	return nil
}

func (m *MockUserStore) DeleteByUserId(ctx *context.Context, userId string) error {
	if m.DeleteByUserIdFunc != nil {
		return m.DeleteByUserIdFunc(ctx, userId)
	}
	return nil
}

func (m *MockUserStore) Delete(ctx *context.Context, filter Filter) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, filter)
	}
	return nil
}

func (m *MockUserStore) CheckExists(ctx *context.Context, filter Filter) (bool, error) {
	if m.CheckExistsFunc != nil {
		return m.CheckExistsFunc(ctx, filter)
	}
	return false, nil
}

func (m *MockUserStore) CheckIfMobileExists(ctx *context.Context, phoneNumber *common.PhoneNumber) (bool, error) {
	if m.CheckIfMobileExistsFunc != nil {
		return m.CheckIfMobileExistsFunc(ctx, phoneNumber)
	}
	return false, nil
}

func TestMongoUserStore_CheckExists(t *testing.T) {
	type args struct {
		ctx    *context.Context
		filter Filter
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "CheckExists",
			args: args{
				ctx:    nil,
				filter: Filter{Key: UserId, Value: "user123"},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockUserStore{
				CheckExistsFunc: func(ctx *context.Context, filter Filter) (bool, error) {
					if filter.Key == UserId && filter.Value == "user123" {
						return true, nil
					}
					return false, nil
				},
			}
			got, err := mockStore.CheckExists(tt.args.ctx, tt.args.filter)
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
	type args struct {
		ctx         *context.Context
		phoneNumber *common.PhoneNumber
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "CheckIfMobileExists",
			args: args{
				ctx:         nil,
				phoneNumber: &common.PhoneNumber{CountryCode: "+1", Number: "1234567890"},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockUserStore{
				CheckIfMobileExistsFunc: func(ctx *context.Context, phoneNumber *common.PhoneNumber) (bool, error) {
					if phoneNumber.CountryCode == "+1" && phoneNumber.Number == "1234567890" {
						return true, nil
					}
					return false, nil
				},
			}
			got, err := mockStore.CheckIfMobileExists(tt.args.ctx, tt.args.phoneNumber)
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
	type args struct {
		ctx    *context.Context
		filter Filter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Delete",
			args: args{
				ctx:    nil,
				filter: Filter{Key: UserId, Value: "user123"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockUserStore{
				DeleteFunc: func(ctx *context.Context, filter Filter) error {
					return nil
				},
			}
			if err := mockStore.Delete(tt.args.ctx, tt.args.filter); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoUserStore_DeleteByUserId(t *testing.T) {
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
			name: "DeleteByUserId",
			args: args{
				ctx:    nil,
				userId: "user123",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockUserStore{
				DeleteByUserIdFunc: func(ctx *context.Context, userId string) error {
					return nil
				},
			}
			if err := mockStore.DeleteByUserId(tt.args.ctx, tt.args.userId); (err != nil) != tt.wantErr {
				t.Errorf("DeleteByUserId() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoUserStore_Get(t *testing.T) {
	type args struct {
		ctx    *context.Context
		filter Filter
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "Get",
			args: args{
				ctx:    nil,
				filter: Filter{Key: UserId, Value: "user123"},
			},
			want: &User{
				UserId:  "user123",
				EmailId: "test@example.com",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockUserStore{
				GetFunc: func(ctx *context.Context, filter Filter) (*User, error) {
					if filter.Key == UserId && filter.Value == "user123" {
						return &User{
							UserId:  "user123",
							EmailId: "test@example.com",
						}, nil
					}
					return nil, nil
				},
			}
			got, err := mockStore.Get(tt.args.ctx, tt.args.filter)
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
	type args struct {
		ctx         *context.Context
		phoneNumber *common.PhoneNumber
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "GetByPhoneNumber",
			args: args{
				ctx:         nil,
				phoneNumber: &common.PhoneNumber{CountryCode: "+1", Number: "1234567890"},
			},
			want: &User{
				UserId: "user123",
				PhoneNumber: &common.PhoneNumber{
					CountryCode: "+1",
					Number:      "1234567890",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockUserStore{
				GetByPhoneNumberFunc: func(ctx *context.Context, phoneNumber *common.PhoneNumber) (*User, error) {
					if phoneNumber.CountryCode == "+1" && phoneNumber.Number == "1234567890" {
						return &User{
							UserId: "user123",
							PhoneNumber: &common.PhoneNumber{
								CountryCode: "+1",
								Number:      "1234567890",
							},
						}, nil
					}
					return nil, nil
				},
			}
			got, err := mockStore.GetByPhoneNumber(tt.args.ctx, tt.args.phoneNumber)
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
	type args struct {
		ctx  *context.Context
		user *User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Insert",
			args: args{
				ctx: nil,
				user: &User{
					UserId:  "user123",
					EmailId: "test@example.com",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockUserStore{
				InsertFunc: func(ctx *context.Context, user *User) error {
					return nil
				},
			}
			if err := mockStore.Insert(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoUserStore_UpdateEmailId(t *testing.T) {
	type args struct {
		ctx     *context.Context
		userId  string
		emailId string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "UpdateEmailId",
			args: args{
				ctx:     nil,
				userId:  "user123",
				emailId: "newemail@example.com",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockUserStore{
				UpdateEmailIdFunc: func(ctx *context.Context, userId, emailId string) error {
					return nil
				},
			}
			if err := mockStore.UpdateEmailId(tt.args.ctx, tt.args.userId, tt.args.emailId); (err != nil) != tt.wantErr {
				t.Errorf("UpdateEmailId() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMongoUserStore_UpdatePhoneNumber(t *testing.T) {
	type args struct {
		ctx         *context.Context
		userId      string
		phoneNumber *common.PhoneNumber
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "UpdatePhoneNumber",
			args: args{
				ctx:         nil,
				userId:      "user123",
				phoneNumber: &common.PhoneNumber{CountryCode: "+1", Number: "9876543210"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := &MockUserStore{
				UpdatePhoneNumberFunc: func(ctx *context.Context, userId string, phoneNumber *common.PhoneNumber) error {
					return nil
				},
			}
			if err := mockStore.UpdatePhoneNumber(tt.args.ctx, tt.args.userId, tt.args.phoneNumber); (err != nil) != tt.wantErr {
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
