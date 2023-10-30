package utils

import (
	"restoran/config"

	"github.com/sashabaranov/go-openai"
)

func OpenAIClient(config config.Config) *openai.Client {
	var client = openai.NewClient(config.AI_API_KEY)

	return client
}
