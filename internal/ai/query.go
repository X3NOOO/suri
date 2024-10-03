package ai

import (
	"context"
	"fmt"

	"github.com/henomis/lingoose/thread"
)

var (
	ErrNoHistory     = fmt.Errorf("no history")
	ErrNoLLMResponse = fmt.Errorf("no LLM response")
)

func (a *SuriAI) Query(query string) (string, error) {
	a.history.Add(thread.NewUserMessage().AddContent(thread.NewTextContent(query)))

	err := a.assistant.RunWithThread(context.Background(), a.history.Get())
	if err != nil {
		return "", err
	}

	h := a.history.Get()

	if len(h.Messages) <= 0 {
		return "", ErrNoHistory
	}

	last_msg := h.Messages[len(h.Messages)-1]

	if last_msg.Role != thread.RoleAssistant {
		return "", ErrNoLLMResponse
	}

	return last_msg.Contents[len(last_msg.Contents)-1].AsString(), nil
}
