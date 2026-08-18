package helper

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWTWithClaims(t *testing.T) {
	signKey := "testSecretKey123"
	jwtService := NewJWT(signKey)

	adminID := "admin-uuid-123"
	adminName := "Bagus Ario Yudanto"
	adminEmail := "bagus@example.com"

	res := jwtService.GenerateJWT(adminID, adminName, adminEmail)
	assert.NotNil(t, res)
	tokenString, ok := res["access_token"].(string)
	assert.True(t, ok)
	assert.NotEmpty(t, tokenString)

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	})
	assert.Nil(t, err)
	assert.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, adminID, claims["id"])
	assert.Equal(t, adminName, claims["name"])
	assert.Equal(t, adminEmail, claims["email"])
	assert.NotNil(t, claims["exp"])
	assert.NotNil(t, claims["iat"])
}
