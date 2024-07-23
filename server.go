package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/X3NOOO/suri/internal/ai"
	"github.com/X3NOOO/suri/routes"
)

type suri struct {
	addr     string
	ai       ai.AI // this value is later passed down to the routing context, it's not used by this struct or any of it's methods directly
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
func (s *suri) WithAI(ai ai.AI) *suri {
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

	ctx := &routes.RoutingContext{
		AI:        s.ai,
		MaxMemory: maxMemory << 20, // convert MB to bytes
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
	router.Handle("GET /knowladge", http.HandlerFunc(ctx.KnowladgeGET))                  //  \VVVV/    this one as well
	router.Handle("POST /knowladge", http.HandlerFunc(ctx.KnowladgePOST))                //   \V/
	router.Handle("GET /knowladge/{name}", http.HandlerFunc(ctx.KnowladgeNameGET))       // {name} endpoints are not implemented yet as we still use
	router.Handle("PUT /knowladge/{name}", http.HandlerFunc(ctx.KnowladgeNamePUT))       // the json database
	router.Handle("DELETE /knowladge/{name}", http.HandlerFunc(ctx.KnowladgeNameDELETE)) //  - all return 501

	log.Println("Starting Suri server on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
