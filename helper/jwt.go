package helper

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTInterface interface {
	GenerateJWT(userID string) map[string]any
}

type JWT struct {
	signKey string
}

func New(signKey string) JWTInterface {
	return &JWT{
		signKey: signKey,
	}
}

func (j *JWT) GenerateJWT(adminID string) map[string]any {
	var result = map[string]any{}
	var accessToken = j.generateToken(adminID)
	if accessToken == "" {
		return nil
	}
	result["access_token"] = accessToken
	return result
}

func (j *JWT) generateToken(id string) string {
	var claims = jwt.MapClaims{}
	claims["id"] = id
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()

	var sign = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validToken, err := sign.SignedString([]byte(j.signKey))

	if err != nil {
		return ""
	}

	return validToken
}
