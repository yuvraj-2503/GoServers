package token_manager

import (
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
)

func (j *JwtTokenManager) Generate(claims *TokenClaims) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	setClaims(claims, token)
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}

	// Base64URL encoding the generated token
	return base64.RawURLEncoding.EncodeToString([]byte(tokenString)), nil
}

func setClaims(tokenClaims *TokenClaims, token *jwt.Token) {
	claims := token.Claims.(jwt.MapClaims)
	claims["iat"] = tokenClaims.IAT.Unix()
	claims["exp"] = tokenClaims.EXP.Unix()
	claims["authorized"] = true
	claims["user"] = tokenClaims.UserId
	claims["email"] = tokenClaims.EmailId
	claims["machine"] = tokenClaims.MachineId
	claims["app"] = tokenClaims.App
	claims["kind"] = tokenClaims.Kind
	// Added for compatibility
	claims["machine_id"] = tokenClaims.MachineId
	claims["jti"] = tokenClaims.App
	claims["sub"] = tokenClaims.Sub
	claims["userId"] = tokenClaims.UserId
}
