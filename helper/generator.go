package helper

import "github.com/google/uuid"

func GenerateUUID() string {
	var result = uuid.NewString()

	return result
}
