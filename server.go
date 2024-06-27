package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/X3NOOO/suri/internal/ai"
	"github.com/X3NOOO/suri/routes"
)

type suri struct {
	addr     string
	ai       ai.AI // this value is later passed down to the routing context, it's not used by this struct or any of it's methods per se
	metadata map[string]any
}

/*
Create a new Suri server instance

Args:

	addr: address to bind the server to
	metadata: metadata to be returned by the server (as json, on GET /)
*/
func NewSuri(addr string, metadata map[string]any, historyCountdown time.Duration) *suri {
	return &suri{
		addr:     addr,
		ai:       ai.New(historyCountdown),
		metadata: metadata,
	}
}

// Use custom AI implementation
func (s *suri) WithAI(ai ai.AI) *suri {
	s.ai = ai
	return s
}

func (s *suri) Start() error {
	ctx := &routes.RoutingContext{
		AI: s.ai,
	}

	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		msg, err := json.Marshal(s.metadata)
		if err != nil {
			log.Printf("Error while marshaling server metadata: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(msg))
	})

	router.Handle("/status", http.HandlerFunc(ctx.StatusALL))
	router.Handle("POST /query", http.HandlerFunc(ctx.QueryPOST))

	log.Println("Starting Suri server on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
