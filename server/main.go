package main

import (
	"os"
	"strconv"

	"github.com/maelvls/quote/server/service"
	log "github.com/sirupsen/logrus"
)

// Set during build, e.g.: go build  -ldflags"-X main.version=$(git describe)".
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// I'll try to use the go-idiomatic 'happy-path' way.
func main() {

	// Set the port accoding to PORT.
	port := 8000
	if str := os.Getenv("PORT"); str != "" {
		var err error
		port, err = strconv.Atoi(str)
		if err != nil {
			log.Fatalf("PORT not valid: %v", str)
		}
	}

	// Set the log format according to LOG_FORMAT.
	switch os.Getenv("LOG_FORMAT") {
	case "":
		log.SetFormatter(&log.TextFormatter{})
	case "json":
		log.SetFormatter(&log.JSONFormatter{})
	default:
		log.Fatal(`Expected LOG_FORMAT is 'json'`)
	}

	log.Printf("serving on port %v (version %s)", port, version)

	s := service.NewServer()
	if err := s.Run(port); err != nil {
		log.Fatalf("launching server: %v", err)
	}
}
