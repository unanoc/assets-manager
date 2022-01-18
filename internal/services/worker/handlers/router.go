package handlers

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/trustwallet/assets-manager/internal/services/worker/events"
)

func NewRouter(eh *events.EventHandler) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/", GithubEventsHandler(eh))

	return mux
}
