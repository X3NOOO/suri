package ai

import "strconv"

const (
	defaultModel           = "gpt-3.5-turbo"
	defaultEmbeddingModel  = "text-embedding-3-small"
	defaultTemperatureStr  = "0.5"
	defaultSystemMessage   = "You're a kind and polite personal assistant."
	defaultRagPath         = "rag.json"
	defaultCachePath       = "cache.json"
	defaultChunkSizeStr    = "1000"
	defaultChunkOverlapStr = "50"
)

var defaultTemperature, _ = strconv.ParseFloat(defaultTemperatureStr, 32)
var defaultChunkSize, _ = strconv.Atoi(defaultChunkSizeStr)
var defaultChunkOverlap, _ = strconv.Atoi(defaultChunkOverlapStr)
