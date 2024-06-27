package ai

import "strconv"

const (
	defaultModel          = "gpt-3.5-turbo"
	defaultEmbeddingModel = "text-embedding-3-small"
	defaultTemperatureStr = "0.5"
	defaultSystemMessage  = "You're kind and polite personal assistant."
)

var defaultTemperature, _ = strconv.ParseFloat(defaultTemperatureStr, 32)
