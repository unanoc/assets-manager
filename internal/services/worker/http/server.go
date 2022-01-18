package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/trustwallet/assets-manager/internal/config"
	"github.com/trustwallet/assets-manager/internal/services/worker/events"
)

type Server struct {
	srv *http.Server
}

// Run runs a server.
func (s Server) Run() {
	log.Infof("Running HTTP server at :%d", config.Default.Port)

	go func() {
		if err := s.srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

// Shutdown stops a server.
func (s Server) Shutdown(ctx context.Context) error {
	log.Info("Stopping HTTP server")

	return s.srv.Shutdown(ctx)
}

// NewHTTPServer returns an instance of Server with registered HTTP router.
func NewHTTPServer(eh *events.EventHandler) *Server {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", GithubEventsHandler(eh))

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", config.Default.Port),
	}

	return &Server{srv: srv}
}
