package helper

import "github.com/google/uuid"

type GeneratorInterface interface {
	GenerateUUID() (string, error)
}

type Generator struct{}

func NewGenerator() GeneratorInterface {
	return &Generator{}
}

func (g Generator) GenerateUUID() (string, error) {
	result, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return result.String(), nil
}
