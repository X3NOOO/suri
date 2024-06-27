package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var version string = "dev"

var envFlag string
var verboseFlag bool

func main() {
	flag.StringVar(&envFlag, "env", ".env", "Environment file")
	flag.BoolVar(&verboseFlag, "verbose", false, "Verbose mode")
	flag.Parse()

	var logfile io.Writer
	if verboseFlag {
		logfile = os.Stderr
	} else {
		logfile = io.Discard
	}

	log.SetOutput(logfile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	_ = godotenv.Load(envFlag)
	// log.Printf("Environ: %v\n", os.Environ())

	port, ok := os.LookupEnv("SURI_PORT")
	if !ok {
		port = defaultPort
		log.Printf("SURI_PORT environment variable not found. Using default (%s)\n", defaultPort)
	}

	metadata := map[string]any{
		"name":        "Suri",
		"description": "Home AI Assistant (server)",
		"version":     version,
		"source":      "https://github.com/X3NOOO/suri",
	}

	historyCountdownStr, ok := os.LookupEnv("SURI_HISTORY_COUNTDOWN")
	if !ok {
		log.Printf("SURI_HISTORY_COUNTDOWN environment variable not found. Using default (%s)\n", defaultHistoryCountdown)
		historyCountdownStr = defaultHistoryCountdown
	}

	historyCountdown, err := strconv.Atoi(historyCountdownStr)
	if err != nil {
		log.Fatalf("Error while parsing SURI_HISTORY_COUNTDOWN environment variable: %v\n", err)
	}

	server := NewSuri(
		":"+port,
		metadata,
		time.Duration(historyCountdown)*time.Second,
	)

	_ = server

	log.Fatalln(server.Start())
}
