package worker

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	assetsmanager "github.com/trustwallet/assets-go-libs/client/assets-manager"
	"github.com/trustwallet/assets-manager/internal/config"
	"github.com/trustwallet/assets-manager/internal/queue"
	"github.com/trustwallet/assets-manager/internal/services"
	"github.com/trustwallet/assets-manager/internal/services/worker/blockchain"
	"github.com/trustwallet/assets-manager/internal/services/worker/events"
	"github.com/trustwallet/assets-manager/internal/services/worker/github"
	"github.com/trustwallet/assets-manager/internal/services/worker/metrics"
	metricsLib "github.com/trustwallet/go-libs/metrics"
	"github.com/trustwallet/go-libs/mq"
	"github.com/trustwallet/go-libs/worker"
)

type App struct {
	mqClient      *mq.Client
	eventHandler  *events.Handler
	metricsPusher worker.Worker
}

func NewApp() *App {
	services.Setup()

	githubClient, err := github.NewClient()
	if err != nil {
		log.WithError(err).Fatal("failed to create github instance")
	}

	mqClient, err := mq.Connect(config.Default.Rabbitmq.URL)
	if err != nil {
		log.WithError(err).Fatal("failed to init RabbitMQ")
	}

	metricsPusher, err := metrics.InitMetricsPusher(config.Default.Metrics.PushGatewayURL,
		config.Default.Metrics.PushInterval)
	if err != nil {
		log.WithError(err).Error("failed to init metrics pusher")
	}

	assetsManagerClient := assetsmanager.InitClient(config.Default.Clients.AssetsManager.API, nil)
	blockchainClient := blockchain.NewClient()
	prometheus := metrics.NewPrometheus()
	eventHandler := events.NewHandler(prometheus, githubClient, blockchainClient, &assetsManagerClient)

	return &App{
		mqClient:      mqClient,
		eventHandler:  eventHandler,
		metricsPusher: metricsPusher,
	}
}

func (a *App) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	a.mqClient.ListenConnectionAsync(ctx, wg)
	runBackgroundChecker(ctx, wg, a.eventHandler)

	err := a.mqClient.StartConsumers(ctx, initConsumers(ctx, a.mqClient, a.eventHandler)...)
	if err != nil {
		log.WithError(err).Fatal("failed to start Rabbit MQ consumers")
	}

	if a.metricsPusher != nil {
		a.metricsPusher.Start(ctx, wg)
	}

	<-stop

	cancel()
	wg.Wait()
}

func runBackgroundChecker(ctx context.Context, wg *sync.WaitGroup, eh *events.Handler) {
	repoOwner := config.Default.Github.RepoOwner
	repoName := config.Default.Github.RepoName

	workerOpts := worker.DefaultWorkerOptions(config.Default.Timeout.BackgroundCheck)
	workerFn := func() error {
		return eh.CheckOpenPullRequests(ctx, repoOwner, repoName, nil)
	}

	worker.InitWorker("PR checker", workerOpts, workerFn).Start(ctx, wg)
}

func initConsumers(ctx context.Context, mqClient *mq.Client, eh *events.Handler) []mq.Consumer {
	options := mq.DefaultConsumerOptions(config.Default.Consumer.Workers)

	options.PerformanceMetric = metricsLib.NewPerformanceMetric(
		"assets_manager_worker",
		prometheus.Labels{"queue_name": string(queue.QueueAssetManagerProcessGithubEvent)},
		prometheus.DefaultRegisterer,
	)

	consumers := []mq.Consumer{
		mqClient.InitConsumer(queue.QueueAssetManagerProcessGithubEvent, options, events.GetEventConsumer(ctx, eh)),
	}

	return consumers
}
