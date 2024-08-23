package token_manager

import (
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func (j *JwtTokenManager) Verify(apiToken string) (*TokenClaims, error) {
	decodedToken, decodeErr := base64.RawURLEncoding.DecodeString(apiToken)
	if decodeErr != nil {
		return nil, &TokenDecodeError{}
	}
	token, err := jwt.Parse(string(decodedToken), func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})
	if err != nil {
		return nil, &InvalidTokenError{}
	}
	tokenClaims := &TokenClaims{}
	if token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		setJavaClaims(claims, tokenClaims)
		setGoClaims(claims, tokenClaims)
		setCommonClaims(claims, tokenClaims)
		return tokenClaims, nil
	} else {
		return nil, &TokenExpiryError{}
	}
}

// Java claims for compatibility , should be removed in the future
func setJavaClaims(jwtClaims jwt.MapClaims, tokenClaims *TokenClaims) {
	machineId, ok := jwtClaims["machine_id"].(string)
	if ok {
		tokenClaims.MachineId = machineId
	}
	jti, ok := jwtClaims["jti"].(string)
	if ok {
		tokenClaims.JTI = jti
	}
	sub, ok := jwtClaims["sub"].(string)
	if ok {
		tokenClaims.Sub = sub
	}
	userId, ok := jwtClaims["userId"].(string)
	if ok {
		tokenClaims.UserId = userId
	}
}

func setGoClaims(claims jwt.MapClaims, tokenClaims *TokenClaims) {
	userId, ok := claims["user"].(string)
	if ok {
		tokenClaims.UserId = userId
	}
	emailId, ok := claims["email"].(string)
	if ok {
		tokenClaims.EmailId = emailId
	}
	machineId, ok := claims["machine"].(string)
	if ok {
		tokenClaims.MachineId = machineId
	}
	app, ok := claims["app"].(string)
	if ok {
		tokenClaims.App = app
	}
	kind, ok := claims["kind"].(string)
	if ok {
		tokenClaims.Kind = kind
	}
}

func setCommonClaims(claims jwt.MapClaims, tokenClaims *TokenClaims) {
	tokenClaims.IAT = getIATClaim(claims)
	tokenClaims.EXP = getEXPClaim(claims)
}

func getIATClaim(claims jwt.MapClaims) *time.Time {
	iatClaim, ok := claims["iat"].(float64)
	if ok {
		var iatTime = time.Unix(int64(iatClaim), 0)
		return &iatTime
	}
	return nil
}

func getEXPClaim(claims jwt.MapClaims) *time.Time {
	expClaim, ok := claims["exp"].(float64)
	if ok {
		var expTime = time.Unix(int64(expClaim), 0)
		return &expTime
	}
	return nil
}
