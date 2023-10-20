package helper

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTInterface interface {
	GenerateJWT(userID string) map[string]any
	GenerateTableToken(noTable int, adminName string) string
	ExtractToken(tokenString string) (int, error)
}

type JWT struct {
	signKey string
}

func NewJWT(signKey string) JWTInterface {
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
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	var sign = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validToken, err := sign.SignedString([]byte(j.signKey))

	if err != nil {
		return ""
	}

	return validToken
}

func (j *JWT) GenerateTableToken(noTable int, adminName string) string {
	var claims = jwt.MapClaims{}
	claims["no_table"] = noTable
	claims["admin"] = adminName
	claims["iat"] = time.Now().Unix()

	var sign = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	validToken, err := sign.SignedString([]byte(j.signKey))

	if err != nil {
		return ""
	}

	return validToken
}

func (j *JWT) ExtractToken(tokenString string) (int, error) {
	parts := strings.Split(tokenString, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, errors.New("invalid token")
	}
	jwtToken := parts[1]

	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return 0, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		noTableString := fmt.Sprint(claims["no_table"])
		noTable, _ := strconv.Atoi(noTableString)
		return noTable, nil
	}

	return 0, errors.New("claims not found")
}
