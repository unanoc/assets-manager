package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/trustwallet/assets-manager/internal/config"
)

type Server struct {
	srv *http.Server
}

// Run runs a server.
func (s Server) Run() {
	log.WithField("port", config.Default.Port).Infof("Running HTTP server")

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()
}

// Shutdown stops a server.
func (s Server) Shutdown(ctx context.Context) error {
	log.Info("Stopping HTTP server")

	return s.srv.Shutdown(ctx)
}

// NewHTTPServer returns an initialized HTTP server.
func NewHTTPServer(router http.Handler) *Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Default.Port),
		Handler: router,
	}

	return &Server{srv: srv}
}
