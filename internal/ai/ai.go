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
	rag       *rag.RAG
	history   history
}

func New(historyCountdown time.Duration) *SuriAI {
	systemMessage, ok := os.LookupEnv("SURI_SYSTEM_MESSAGE") // LookupEnv instead of Getenv because the systemMessage can be empty I guess
	if !ok {
		systemMessage = defaultSystemMessage
		log.Printf("SURI_SYSTEM_MESSAGE environment variable not found. Using the default (%s)\n", defaultSystemMessage)
	}

	modelStr := os.Getenv("SURI_MODEL")
	if modelStr == "" {
		modelStr = defaultModel
		log.Printf("SURI_MODEL environment variable not found. Using the default (%s)\n", defaultModel)
	}

	model := openai.Model(modelStr)

	temperatureStr := os.Getenv("SURI_TEMPERATURE")
	if temperatureStr == "" {
		temperatureStr = defaultTemperatureStr
		log.Printf("SURI_TEMPERATURE environment variable not found. Using the default (%s)\n", defaultTemperatureStr)
	}

	temperature, err := strconv.ParseFloat(temperatureStr, 32)
	if err != nil {
		temperature = defaultTemperature
		log.Printf("Error while parsing SURI_TEMPERATURE environment variable: %v. Using the default (%s)\n", err, defaultTemperatureStr)
	}

	embeddingModelStr := os.Getenv("SURI_EMBEDDING_MODEL")
	if embeddingModelStr == "" {
		embeddingModelStr = defaultEmbeddingModel
		log.Printf("SURI_EMBEDDING_MODEL environment variable not found. Using the default (%s)\n", defaultEmbeddingModel)
	}

	embeddingModel := openaiembedder.Model(embeddingModelStr)

	ragPath := os.Getenv("SURI_RAG_PATH")
	if ragPath == "" {
		ragPath = defaultRagPath
		log.Printf("SURI_RAG_PATH environment variable not found. Using the default (%s)\n", defaultRagPath)
	}

	cachePath := os.Getenv("SURI_CACHE_PATH")
	if cachePath == "" {
		cachePath = defaultCachePath
		log.Printf("SURI_CACHE_PATH environment variable not found. Using the default (%s)\n", defaultCachePath)
	}

	chunkSizeStr := os.Getenv("SURI_RAG_CHUNK_SIZE")
	if chunkSizeStr == "" {
		chunkSizeStr = defaultChunkSizeStr
		log.Printf("SURI_RAG_CHUNK_SIZE environment variable not found. Using the default (%s)\n", defaultChunkSizeStr)
	}

	chunkSize, err := strconv.Atoi(chunkSizeStr)
	if err != nil {
		chunkSize = defaultChunkSize
		log.Printf("Error while parsing SURI_RAG_CHUNK_SIZE environment variable: %v. Using the default (%d)\n", err, defaultChunkSize)
	}

	chunkOverlapStr := os.Getenv("SURI_RAG_CHUNK_OVERLAP")
	if chunkOverlapStr == "" {
		chunkOverlapStr = defaultChunkOverlapStr
		log.Printf("SURI_RAG_CHUNK_OVERLAP environment variable not found. Using the default (%s)\n", defaultChunkOverlapStr)
	}

	chunkOverlap, err := strconv.Atoi(chunkOverlapStr)
	if err != nil {
		chunkOverlap = defaultChunkOverlap
		log.Printf("Error while parsing SURI_RAG_CHUNK_OVERLAP environment variable: %v. Using the default (%d)\n", err, defaultChunkOverlap)
	}

	localRag := rag.New(index.New(jsondb.New().WithPersist(ragPath), openaiembedder.New(embeddingModel))).WithChunkSize(uint(chunkSize)).WithChunkOverlap(uint(chunkOverlap))

	return &SuriAI{
		assistant: *assistant.New(openai.New().WithTemperature(float32(temperature)).WithModel(model).WithCache(
			cache.New(index.New(jsondb.New().WithPersist(cachePath), openaiembedder.New(embeddingModel))),
		)).WithRAG(localRag),
		rag: localRag,
		history: history{
			Countdown:     historyCountdown,
			SystemMessage: systemMessage,
		},
	}
}
