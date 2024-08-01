package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/X3NOOO/suri/internal/ai"
	"github.com/X3NOOO/suri/internal/tts"
	"github.com/X3NOOO/suri/routes"
	"github.com/X3NOOO/whisper-go"
)

var (
	ErrNoWhisperModel         = errors.New("SURI_STT_WHISPER_MODEL environment variable not found")
	ErrNoPiperPath            = errors.New("PIPER_PATH environment variable not found")
	ErrNoPiperModelPath       = errors.New("PIPER_MODEL_PATH environment variable not found")
	ErrNoPiperModelConfigPath = errors.New("PIPER_MODEL_CONFIG_PATH environment variable not found")
)

type suri struct {
	addr     string
	ai       routes.AI // this value is later passed down to the routing context, it's not used by this struct or any of it's methods directly
	metadata map[string]any
}

/*
Create a new Suri server instance.

Args:

	addr: Address to bind the server to.
	metadata: Metadata to be returned by the server (as json, on GET /).
*/
func NewSuri(addr string, metadata map[string]any, historyCountdown time.Duration) *suri {
	return &suri{
		addr:     addr,
		ai:       ai.New(historyCountdown),
		metadata: metadata,
	}
}

// Use custom AI implementation
func (s *suri) WithAI(ai routes.AI) *suri {
	s.ai = ai
	return s
}

func (s *suri) Start() error {
	maxMemoryStr := os.Getenv("SURI_API_MAX_FILE_SIZE_MB")
	if maxMemoryStr == "" {
		maxMemoryStr = defaultMaxMemoryStr
		log.Printf("SURI_API_MAX_FILE_SIZE_MB environment variable not found. Using the default (%s)\n", defaultMaxMemoryStr)
	}

	maxMemory, err := strconv.ParseInt(maxMemoryStr, 10, 32)
	if err != nil {
		maxMemory = defaultMaxMemory
		log.Printf("Error while parsing SURI_API_MAX_FILE_SIZE_MB environment variable: %v. Using the default (%d)\n", err, defaultMaxMemory)
	}

	whisperModel := os.Getenv("SURI_STT_WHISPER_MODEL")
	if whisperModel == "" {
		return ErrNoWhisperModel
	}
	whisperTempStr := os.Getenv("SURI_STT_WHISPER_TEMPERATURE")
	if whisperTempStr == "" {
		whisperTempStr = defaultWhisperTempStr
	}
	whisperTemp, err := strconv.ParseFloat(whisperTempStr, 64)
	if err != nil {
		whisperTemp = defaultWhisperTemp
		log.Printf("Error while parsing SURI_STT_WHISPER_TEMPERATURE environment variable: %v. Using the default (%f)\n", err, defaultWhisperTemp)
	}
	whisperLanguage := os.Getenv("SURI_STT_LANGUAGE")

	piperPath := os.Getenv("PIPER_BIN_PATH")
	if piperPath == "" {
		return ErrNoPiperPath
	}
	piperModelPath := os.Getenv("PIPER_MODEL_PATH")
	if piperModelPath == "" {
		return ErrNoPiperModelPath
	}
	piperModelConfigPath := os.Getenv("PIPER_MODEL_CONFIG_PATH")
	if piperModelConfigPath == "" {
		return ErrNoPiperModelConfigPath
	}

	ctx_tts := tts.PiperTTS{
		BinPath:         piperPath,
		ModelPath:       piperModelPath,
		ModelConfigPath: piperModelConfigPath,
	}

	ctx := &routes.RoutingContext{
		AI:              s.ai,
		MaxMemory:       maxMemory << 20, // convert MB to bytes
		Whisper:         *whisper.New(os.Getenv("GROQ_API_KEY")),
		WhisperModel:    whisperModel,
		WhisperTemp:     whisperTemp,
		WhisperLanguage: whisperLanguage,
		TTS:             ctx_tts, // tu sie jebie, zla sygnatura
	}

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		msg, err := json.Marshal(s.metadata)
		if err != nil {
			log.Printf("Error while marshaling server metadata: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(msg))
	})

	router.Handle("/status", http.HandlerFunc(ctx.StatusALL))                            //  __\_|__
	router.Handle("POST /query", http.HandlerFunc(ctx.QueryPOST))                        // \VVVVVV/
	router.Handle("GET /knowledge", http.HandlerFunc(ctx.KnowledgeGET))                  //  \VVVV/    this one as well
	router.Handle("POST /knowledge", http.HandlerFunc(ctx.KnowledgePOST))                //   \V/
	router.Handle("GET /knowledge/{name}", http.HandlerFunc(ctx.KnowledgeNameGET))       // {name} endpoints are not implemented yet as we still use
	router.Handle("PUT /knowledge/{name}", http.HandlerFunc(ctx.KnowledgeNamePUT))       // the json database
	router.Handle("DELETE /knowledge/{name}", http.HandlerFunc(ctx.KnowledgeNameDELETE)) //  - all return 501

	log.Println("Starting Suri server on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
