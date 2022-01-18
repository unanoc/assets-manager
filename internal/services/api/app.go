package api

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/trustwallet/assets-manager/internal/http"
	"github.com/trustwallet/assets-manager/internal/services"
	"github.com/trustwallet/assets-manager/internal/services/api/handlers"
)

type App struct {
	server *http.Server
}

func NewApp() *App {
	services.Setup()

	router := handlers.NewRouter()
	server := http.NewHTTPServer(router)

	return &App{
		server: server,
	}
}

func (a *App) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	a.server.Run()
	<-done
	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown http server: %v", err)
	}

	cancel()
	wg.Wait()
}
