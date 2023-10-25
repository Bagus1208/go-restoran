package helper

import "github.com/google/uuid"

type GeneratorInterface interface {
	GenerateUUID() (string, error)
}

type generator struct{}

func NewGenerator() GeneratorInterface {
	return &generator{}
}

func (generate generator) GenerateUUID() (string, error) {
	result, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return result.String(), nil
}
