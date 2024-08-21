package otp

import (
	"context"
	"crypto/sha1"
	"encoding/binary"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"otp-manager/common"
	"otp-manager/senders"
	"time"
)

type OtpManager interface {
	Send(ctx *context.Context, contact *common.Contact) (*string, error)
	Verify(ctx *context.Context, sessionId string, otp uint64) (*common.Contact, error)
}

type MongoOtpManager struct {
	otpStore OTPStore
	sender   senders.OtpSender
}

func (m *MongoOtpManager) Send(ctx *context.Context, contact *common.Contact) (*string, error) {
	err := make(chan error)
	sessionId := sessionId()
	otp := generateOTP()
	m.storeOtp(ctx, sessionId, otp, contact, err)
	m.send(contact, otp, err)
	for i := 0; i < 2; i++ {
		e := <-err
		if e != nil {
			return nil, e
		}
	}
	return &sessionId, nil
}

func (m *MongoOtpManager) Verify(ctx *context.Context, sessionId string, otp uint64) (*common.Contact, error) {
	return nil, nil
}

func (m *MongoOtpManager) storeOtp(ctx *context.Context, sessionId string, otp uint64, contact *common.Contact, err chan error) {
	go func() {
		sha := computeSha1(otp)
		err <- m.otpStore.Upsert(ctx, &common.OTP{
			SessionId: sessionId,
			Contact:   contact,
			Otp:       sha,
			Retries:   0,
			CreatedOn: time.Now(),
		})
	}()
}

func (m *MongoOtpManager) send(contact *common.Contact, otp uint64, err chan error) {
	go func() {
		err <- m.sender.Send(contact, otp)
	}()
}

func computeSha1(otp uint64) []byte {
	hasher := sha1.New()
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, otp)
	hasher.Write(bs)
	hash := hasher.Sum(nil)
	return hash
}

func generateOTP() uint64 {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	min := 100000
	max := 999999
	otp := uint64(rng.Intn(max-min+1) + min)
	return otp
}

func sessionId() string {
	bsonId := primitive.NewObjectID().Hex()
	return bsonId
}
