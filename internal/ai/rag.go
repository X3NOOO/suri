package ai

import (
	"context"

	"github.com/X3NOOO/suri/internal/parser"
	"github.com/henomis/lingoose/document"
)

func (a *SuriAI) Learn(documents ...document.Document) error {
	return a.rag.AddDocuments(context.Background(), documents...)
}

func (a *SuriAI) LearnFile(file []byte, mime string, metadata map[string]any) error {
	content, err := parser.Parse(file, mime) // We use a custom parser package instead of the lingoose's built-in one because it requires the file to be saved on disk
	if err != nil {
		return err
	}

	a.rag.AddDocuments(context.Background(), document.Document{
		Content:  content,
		Metadata: metadata,
	})

	return nil
}
