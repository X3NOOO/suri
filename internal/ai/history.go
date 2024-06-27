package ai

import (
	"github.com/henomis/lingoose/thread"
	"time"
)

type history struct {
	Countdown     time.Duration
	SystemMessage string

	lastAccess time.Time
	history    thread.Thread
}

func (h *history) Add(message *thread.Message) {
	if time.Since(h.lastAccess) > h.Countdown {
		h.history = thread.Thread{}

		if h.SystemMessage != "" {
			h.history.AddMessage(thread.NewSystemMessage().AddContent(thread.NewTextContent(h.SystemMessage)))
		}
	}

	h.history.AddMessage(message)

	h.lastAccess = time.Now()
}

func (h *history) Get() *thread.Thread {
	if time.Since(h.lastAccess) > h.Countdown {
		h.history = thread.Thread{}

		if h.SystemMessage != "" {
			h.history.AddMessage(thread.NewSystemMessage().AddContent(thread.NewTextContent(h.SystemMessage)))
		}
	}

	h.lastAccess = time.Now()

	return &h.history
}
