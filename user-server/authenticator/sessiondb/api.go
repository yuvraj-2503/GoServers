package sessiondb

import (
	"context"
	"time"
)

type SessionMapping struct {
	SessionId string     `bson:"sessionId"`
	UserId    string     `bson:"userId"`
	CreatedOn *time.Time `bson:"createdOn"`
}

type SessionRepository interface {
	Insert(ctx *context.Context, sessionMapping *SessionMapping) error
	Get(ctx *context.Context, sessionId string) (*SessionMapping, error)
	Delete(ctx *context.Context, sessionId string) error
}
