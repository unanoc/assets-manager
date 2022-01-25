package worker

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/trustwallet/assets-manager/internal/config"
	"github.com/trustwallet/assets-manager/internal/http"
	"github.com/trustwallet/assets-manager/internal/services"
	"github.com/trustwallet/assets-manager/internal/services/worker/blockchain"
	"github.com/trustwallet/assets-manager/internal/services/worker/events"
	"github.com/trustwallet/assets-manager/internal/services/worker/github"
	"github.com/trustwallet/assets-manager/internal/services/worker/handlers"
	"github.com/trustwallet/assets-manager/internal/services/worker/metrics"
	assetsmanager "github.com/trustwallet/go-libs/client/api/assets-manager"
	"github.com/trustwallet/go-libs/worker"
)

type App struct {
	metrics      *metrics.Prometheus
	eventHandler *events.EventHandler
	server       *http.Server
}

func NewApp() *App {
	services.Setup()

	githubClient, err := github.NewClient()
	if err != nil {
		log.Fatalf("failed to create github instance: %v", err)
	}

	assetsManagerClient := assetsmanager.InitClient(config.Default.Clients.AssetsManagerAPI, nil)
	blockchainClient := blockchain.NewClient()
	prometheus := metrics.NewPrometheus()
	eventHandler := events.NewEventHandler(prometheus, githubClient, blockchainClient, &assetsManagerClient)

	router := handlers.NewRouter(eventHandler)
	server := http.NewHTTPServer(router)

	return &App{
		metrics:      prometheus,
		eventHandler: eventHandler,
		server:       server,
	}
}

func (a *App) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}

	checkFunc := func() {
		err := a.eventHandler.CheckOpenPullRequests(ctx, config.Default.Github.RepoOwner,
			config.Default.Github.RepoName, nil, false)
		if err != nil {
			log.Error(err)
		}

		err = a.eventHandler.CheckClosedPullRequests(ctx, config.Default.Github.RepoOwner,
			config.Default.Github.RepoName)
		if err != nil {
			log.Error(err)
		}
	}

	checker := worker.New("PR checker", checkFunc, config.Default.Timeout.BackgroundCheck, false)
	checker.StartWithTicker(ctx, wg)

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
