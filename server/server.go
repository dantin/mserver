package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/dantin/mserver/pkg/logutil"
	log "github.com/sirupsen/logrus"
)

// Server implements a HTTP server.
type Server struct {
	server *http.Server

	cfg *Config
}

// NewServer creates a new instance of HTTP server.
func NewServer(cfg *Config) *Server {
	// init logger.
	logutil.InitLogger(&cfg.Log)
	showVersionInfo()

	svr = &Server{
		cfg: cfg,
		server: &http.Server{
			Addr: cfg.ListenAddr,

			ReadTimeout:    5 * time.Second,
			WriteTimeout:   10 * time.Second,
			IdleTimeout:    30 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}

	return svr
}

// Run starts the HTTP server.
func (s *Server) Run() error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	done := make(chan bool)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		sig := <-sc
		fmt.Printf("got signal [%d] to exit.\n", sig)
		log.Infof("%s - Shutdown signal received...", hostname)

		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		s.server.SetKeepAlivesEnabled(false)

		if err := s.server.Shutdown(ctx); err != nil {
			log.Fatalf("Could not gracefully shutdown the server: %v", err)
		}
		close(done)
	}()

	log.Infof("%s - Start server on port %v", hostname, s.server.Addr)
	atomic.StoreInt32(&healthy, 1)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v", s.server.Addr, err)
	}

	<-done
	log.Infof("%s - Server gracefully stopped", hostname)

	return nil
}
