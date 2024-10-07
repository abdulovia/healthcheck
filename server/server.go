package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server interface {
	HandlePing(w http.ResponseWriter, r *http.Request)
}

func NewServer(cfg *Config) Server {
	return &server{
		config: cfg,
	}
}

type server struct {
	config *Config
}

func (s *server) HandlePing(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	log.Printf("[INFO] got healthcheck request from %q", r.RemoteAddr)
	status := 0
	defer func() {
		log.Printf("[INFO] finished healthcheck, status = %v, duration = %s", status, time.Since(start))
	}()

	log.Println("[INFO] got token from metadata, trying to ping all endpoints")

	p := pinger{
		endpoints: struct {
			firEndpoint       string
			secEndpoint         string
		}{
			firEndpoint:       s.config.FirstEndpoint,
			secEndpoint:         s.config.SecondEndpoint,
		},
	}
	err := p.ping(r.Context())
	if err != nil {
		status = 1
		log.Printf("[ERROR] failed to ping services: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("[INFO] successfully finished healthcheck")
}

type pinger struct {
	endpoints struct {
		firEndpoint       string
		secEndpoint         string
	}
}

func (p *pinger) pingFirst(ctx context.Context) error {
	ep := fmt.Sprintf("http://%s/ping", p.endpoints.firEndpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return fmt.Errorf("failed to create request for first end: %s", err)
	}

	cl := &http.Client{
		Timeout: 200 * time.Millisecond,
	}
	resp, err := cl.Do(req)
	if err != nil {
		return fmt.Errorf("failed to ping first end: %s", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to ping first end: got non-OK status %d", resp.StatusCode)
	}
	return nil
}

func (p *pinger) pingSecond(ctx context.Context) error {
	ep := fmt.Sprintf("http://%s/ping", p.endpoints.secEndpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return fmt.Errorf("failed to create request for second end: %s", err)
	}

	cl := &http.Client{
		Timeout: 200 * time.Millisecond,
	}
	resp, err := cl.Do(req)
	if err != nil {
		return fmt.Errorf("failed to ping second end: %s", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to ping second end: %s", resp.Status)
	}
	return nil
}

func (p *pinger) ping(ctx context.Context) error {
	errs := make(chan error, 2)
	go func() {
		errs <- p.pingFirst(ctx)
	}()
	go func() {
		errs <- p.pingSecond(ctx)
	}()
	var retError error
	for i := 0; i < 2; i++ {
		v := <-errs
		if v != nil {
			retError = v
		}
	}
	return retError
}
