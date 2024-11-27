package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	token "token-manager"
	"user-server/common"
)

type AuthHandler struct {
	tokenManager token.TokenManager
}

func NewAuthHandler(secretKey string) *AuthHandler {
	return &AuthHandler{
		tokenManager: token.NewJwtTokenManager(secretKey),
	}
}

func (a *AuthHandler) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			common.Unauthorized(c, "missing-header", "Auth header not found")
			return
		}

		extractedToken := a.extractToken(authHeader)
		if extractedToken == "" {
			common.Unauthorized(c, "missing-token", "invalid or missing Bearer token")
			return
		}

		jwtToken, err := a.validate(extractedToken)
		if err != nil {
			if errors.Is(err, &token.InvalidTokenError{}) || errors.Is(err, &token.TokenDecodeError{}) {
				common.Unauthorized(c, "invalid-token", "Auth token is not valid")
				return
			}
		} else {
			c.Set("user", *jwtToken)
			c.Next()
		}

	}
}

func (a *AuthHandler) validate(token string) (*token.TokenClaims, error) {
	jwtToken, err := a.tokenManager.Verify(token)
	return jwtToken, err
}

func (a *AuthHandler) extractToken(authHeader string) string {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	token := parts[1]
	return token
}
