package token_manager

import (
	"encoding/base64"
	"time"
)

type TokenClaims struct {
	UserId    string
	EmailId   string
	MachineId string
	App       string
	IAT       *time.Time
	EXP       *time.Time
	Kind      string
	Sub       string
	JTI       string
}

type TokenManager interface {
	Generate(claims *TokenClaims) (string, error)
	TokenVerifier
}

type TokenVerifier interface {
	Verify(token string) (*TokenClaims, error)
}

type JwtTokenManager struct {
	secretKey []byte
}

func NewJwtTokenManager(secretKeyBase64 string) *JwtTokenManager {
	keyBytes, _ := base64.StdEncoding.DecodeString(secretKeyBase64)
	return &JwtTokenManager{
		secretKey: keyBytes,
	}
}
