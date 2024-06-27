package ai

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/henomis/lingoose/assistant"
	openaiembedder "github.com/henomis/lingoose/embedder/openai"
	"github.com/henomis/lingoose/index"
	"github.com/henomis/lingoose/index/vectordb/jsondb"
	"github.com/henomis/lingoose/llm/cache"
	"github.com/henomis/lingoose/llm/openai"
	"github.com/henomis/lingoose/rag"
)

type SuriAI struct {
	assistant assistant.Assistant

	history history
}

func New(historyCountdown time.Duration) *SuriAI {
	systemMessage, ok := os.LookupEnv("SURI_SYSTEM_MESSAGE")
	if !ok {
		systemMessage = defaultSystemMessage
		log.Printf("SURI_SYSTEM_MESSAGE environment variable not found. Using the default (%s)\n", defaultSystemMessage)
	}

	modelStr, ok := os.LookupEnv("SURI_MODEL")
	if !ok {
		modelStr = defaultModel
		log.Printf("SURI_MODEL environment variable not found. Using the default (%s)\n", defaultModel)
	}

	model := openai.Model(modelStr)

	temperatureStr, ok := os.LookupEnv("SURI_TEMPERATURE")
	if !ok {
		temperatureStr = defaultTemperatureStr
		log.Printf("SURI_TEMPERATURE environment variable not found. Using the default (%s)\n", defaultTemperatureStr)
	}

	temperature, err := strconv.ParseFloat(temperatureStr, 32)
	if err != nil {
		temperature = defaultTemperature
		log.Printf("Error while parsing SURI_TEMPERATURE environment variable: %v. Using the default (%s)\n", err, defaultTemperatureStr)
	}

	embeddingModelStr, ok := os.LookupEnv("SURI_EMBEDDING_MODEL")
	if !ok {
		embeddingModelStr = defaultEmbeddingModel
		log.Printf("SURI_EMBEDDING_MODEL environment variable not found. Using the default (%s)\n", defaultEmbeddingModel)
	}

	embeddingModel := openaiembedder.Model(embeddingModelStr)

	return &SuriAI{
		assistant: *assistant.New(openai.New().WithTemperature(float32(temperature)).WithModel(model).WithCache(
			cache.New(index.New(jsondb.New().WithPersist("cache.json"), openaiembedder.New(embeddingModel))),
		)).WithRAG(rag.New(index.New(jsondb.New().WithPersist("rag.json"), openaiembedder.New(embeddingModel)))),
		history: history{
			Countdown:     historyCountdown,
			SystemMessage: systemMessage,
		},
	}
}
