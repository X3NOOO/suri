package routes

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/X3NOOO/suri/models"
	"github.com/X3NOOO/whisper-go"
)

func (ctx *RoutingContext) queryParseResponse(w http.ResponseWriter, r *http.Request, response string) {
	responseType := strings.ToLower(r.Header.Get("Accept"))
	noResponse := strings.ToLower(r.Header.Get("X-No-Response")) == "true"
	muteServerAudio := strings.ToLower(r.Header.Get("X-Mute-Server-Audio")) == "true"

	if noResponse {
		log.Println("No-Response requested")
		return
	}

	log.Println("LLM response:", response)

	audio, err := ctx.TTS.Generate(response)
	if err != nil {
		log.Println("Error while generating audio:", err)
	}

	if !muteServerAudio {
		go func() {
			err = audio.Play()
			if err != nil {
				log.Println("Error while playing audio:", err)
			}
		}()
	}

	var raw_response []byte

	switch responseType {
	case "text/plain":
		w.Header().Add("Content-Type", "text/plain")
		raw_response = []byte(response)

	case "audio/wav", "audio/wave", "audio/x-wav":
		w.Header().Add("Content-Type", "audio/wav")
		raw_response, err = audio.Wav()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	default:
		w.Header().Add("Content-Type", "application/json")

		response := models.QueryPOSTResponse{
			Response: response,
		}

		response_json, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		raw_response = response_json
	}

	log.Println("Response:", string(raw_response))

	if raw_response != nil {
		_, err := w.Write(raw_response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (ctx *RoutingContext) queryPOSTPlainText(w http.ResponseWriter, r *http.Request) {
	text, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	llm_response, err := ctx.AI.Query(string(text))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.queryParseResponse(w, r, llm_response)
}

func (ctx *RoutingContext) queryPOSTJson(w http.ResponseWriter, r *http.Request) {
	body_json, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var body models.QueryPOSTRequest

	err = json.Unmarshal(body_json, &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	llm_response, err := ctx.AI.Query(body.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.queryParseResponse(w, r, llm_response)
}

func (ctx *RoutingContext) queryPOSTAudio(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(ctx.MaxMemory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// FIXME: automatically convert the content type to one that's supported

	query, err := ctx.Whisper.Transcribe(whisper.Request{
		File: whisper.File{
			Data: fileContent,
			Name: "file.wav",
		},
		Model:          ctx.WhisperModel,
		Temperature:    ctx.WhisperTemp,
		ResponseFormat: "text",
		Language:       ctx.WhisperLanguage,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	text, ok := (*query)["text"].(string)
	if !ok {
		http.Error(w, "Failed to decode the transcription", http.StatusInternalServerError)
		return
	}

	log.Println("Audio transcription:", text)

	llm_response, err := ctx.AI.Query(text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx.queryParseResponse(w, r, llm_response)
}

func (ctx *RoutingContext) QueryPOST(w http.ResponseWriter, r *http.Request) {
	noResponse := strings.ToLower(r.Header.Get("X-No-Response")) == "true"
	muteServerAudio := strings.ToLower(r.Header.Get("X-Mute-Server-Audio")) == "true"

	if noResponse && muteServerAudio {
		http.Error(w, "Both X-No-Response and X-Mute-Server-Audio headers are set to true. NOOP requested.", http.StatusBadRequest)
		return
	}

	contentType := strings.Split(strings.ToLower(r.Header.Get("Content-Type")), ";")[0]

	switch contentType { // I was thinking about putting this in a [string]func(){...}map
	case "text/plain":
		ctx.queryPOSTPlainText(w, r)
	case "application/json":
		ctx.queryPOSTJson(w, r)
	case "multipart/form-data":
		ctx.queryPOSTAudio(w, r)
	default:
		http.Error(w, "Unsupported content type", http.StatusBadRequest)
	}
}
