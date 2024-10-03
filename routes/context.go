package routes

import (
	"github.com/X3NOOO/whisper-go"
	"github.com/henomis/lingoose/document"
)

type AI interface {
	Query(query string) (string, error)
	Learn(...document.Document) error
	LearnFile(file []byte, mime string, metadata map[string]any) error
}

type Audio interface {
	Play() error
	Wav() ([]byte, error)
}

type TTS interface {
	Generate(text string) (Audio, error)
}

type RoutingContext struct {
	AI        AI
	MaxMemory int64
	Whisper   whisper.Whisper // in the future we're gonna add more Speech-to-Text engines

	WhisperTemp     float64
	WhisperModel    string
	WhisperLanguage string

	TTS TTS
}
