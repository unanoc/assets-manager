package mergefeebot

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/trustwallet/assets-manager/internal/merge-fee-bot/blockchain"
	"github.com/trustwallet/assets-manager/internal/merge-fee-bot/config"
	"github.com/trustwallet/assets-manager/internal/merge-fee-bot/events"
	"github.com/trustwallet/assets-manager/internal/merge-fee-bot/github"
	"github.com/trustwallet/assets-manager/internal/merge-fee-bot/http"
	"github.com/trustwallet/assets-manager/internal/merge-fee-bot/metrics"
	"github.com/trustwallet/go-libs/client/api/backend"
	"github.com/trustwallet/go-libs/worker"
)

type App struct {
	metrics      *metrics.Prometheus
	eventHandler *events.EventHandler
	server       *http.Server
}

func NewApp() *App {
	setup()

	githubClient, err := github.NewClient()
	if err != nil {
		log.Fatalf("failed to create github instance: %v", err)
	}

	backendClient := backend.InitClient(config.Default.ClientsURLs.BackendAPI, nil)
	blockchainClient := blockchain.NewClient()
	prometheus := metrics.NewPrometheus()
	eventHandler := events.NewEventHandler(prometheus, githubClient, blockchainClient, &backendClient)
	server := http.NewHTTPServer(eventHandler)

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
		err := a.eventHandler.CheckStatusOfOpenPullRequests(ctx, config.Default.Github.RepoOwner,
			config.Default.Github.RepoName, nil, false)
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

func setup() {
	config.SetConfig()

	logLevel, err := log.ParseLevel(config.Default.App.LogLevel)
	if err != nil {
		log.WithError(err).Fatal("failed to parse log level")
	}

	log.SetLevel(logLevel)
}
