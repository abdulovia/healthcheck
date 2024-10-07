package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/abdulovia/healthcheck/server"
)

var (
	port = flag.String("port", ":8792", "Healthcheck application port")
)

func main() {
	cfg := server.NewConfig()
	flag.Parse()

	http.DefaultClient.Timeout = 300 * time.Millisecond

	s := server.NewServer(cfg())
	http.HandleFunc("/ping", s.HandlePing)
	err := http.ListenAndServe(*port, nil)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to start http healthcheck server: %s", err)
	}
}
