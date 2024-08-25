package otp

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"math/rand"
	"otp-manager/common"
	otpErrors "otp-manager/errors"
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

func NewMongoOtpManager(otpStore OTPStore, sender senders.OtpSender) OtpManager {
	return &MongoOtpManager{
		otpStore: otpStore,
		sender:   sender,
	}
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
	result, err := m.otpStore.Get(ctx, sessionId)
	err = handleOtpErrors(result, err)
	if err != nil {
		return nil, err
	}
	err = verifyOtp(m.otpStore, sessionId, otp, result)
	if err != nil {
		return nil, err
	}
	return result.Contact, nil
}

func verifyOtp(otpStore OTPStore, sessionId string, otp uint64, result *common.OTP) error {
	provided := computeSha1(otp)
	equal := bytes.Equal(provided, result.Otp)
	if equal {
		deleteOtp(otpStore, sessionId)
		return nil
	}
	updateRetryCount(otpStore, sessionId)
	return otpErrors.NewOtpError(otpErrors.INCORRECT, "Invalid OTP")
}

func updateRetryCount(store OTPStore, sessionId string) {
	go func() {
		ctx := context.Background()
		err := store.UpdateRetryCount(&ctx, sessionId)
		if err != nil {
			log.Printf("Error updating otp retry count for sessionId: %s, reason: %v\n", sessionId, err)
		}
	}()
}

func deleteOtp(store OTPStore, sessionId string) {
	go func() {
		ctx := context.Background()
		err := store.Delete(&ctx, sessionId)
		if err != nil {
			log.Printf("Error deleting otp for sessionId: %s, reason: %v\n", sessionId, err)
		}
	}()
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

func handleOtpErrors(result *common.OTP, err error) error {
	if err != nil {
		if errors.Is(err, &otpErrors.OtpDoesNotExistError{}) {
			return otpErrors.NewOtpError(otpErrors.NOT_FOUND, err.Error())
		}
		return err
	}

	if result.Retries > 5 {
		return otpErrors.NewOtpError(otpErrors.LIMIT_EXCEEDED, "Maximum otp verification limit exceeded")
	}

	expTime := result.CreatedOn.Add(time.Duration(10) * time.Minute)
	if expTime.Before(time.Now()) {
		return otpErrors.NewOtpError(otpErrors.EXPIRED, "otp expired")
	}
	return nil
}
