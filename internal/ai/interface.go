package ai

import (
	"github.com/henomis/lingoose/document"
)

type AI interface {
	Query(query string) (string, error)
	Learn(...document.Document) error
	LearnFile(file []byte, mime string, metadata map[string]any) error
}
