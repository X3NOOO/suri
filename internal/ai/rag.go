package ai

import (
	"context"

	"github.com/henomis/lingoose/document"
)

func (a *SuriAI) Learn(documents ...document.Document) error {
	return a.rag.AddDocuments(context.Background(), documents...)
}

func (a *SuriAI) LearnFile(file []byte, mime string, metadata map[string]any) error {
	content, err := a.parser.Parse(file, mime)
	if err != nil {
		return err
	}

	a.rag.AddDocuments(context.Background(), document.Document{
		Content:  content,
		Metadata: metadata,
	})

	return nil
}
