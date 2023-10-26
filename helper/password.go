package helper

import "golang.org/x/crypto/bcrypt"

type HashInterface interface {
	HashPassword(password string) (string, error)
	CompareHash(password string, hash string) bool
}

type hash struct{}

func NewHash() HashInterface {
	return &hash{}
}

func (h hash) HashPassword(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (h hash) CompareHash(password, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	if err != nil {
		return false
	}

	return true
}
