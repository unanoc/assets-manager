package api

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/trustwallet/assets-manager/internal/config"
	"github.com/trustwallet/assets-manager/internal/services"
	"github.com/trustwallet/assets-manager/internal/services/api/handlers"
	"github.com/trustwallet/go-libs/httplib"
	"github.com/trustwallet/go-libs/mq"
)

type App struct {
	server   httplib.Server
	mqClient *mq.Client
}

func NewApp() *App {
	services.Setup()

	mqClient, err := mq.Connect(config.Default.Rabbitmq.URL)
	if err != nil {
		log.WithError(err).Fatal("failed to init Rabbit MQ")
	}

	router := handlers.NewRouter(mqClient)
	server := httplib.NewHTTPServer(router, strconv.Itoa(config.Default.Port))

	return &App{
		server:   server,
		mqClient: mqClient,
	}
}

func (a *App) Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	wg := &sync.WaitGroup{}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	a.mqClient.ListenConnectionAsync(ctx, wg)
	a.server.Run(ctx, wg)

	<-stop

	cancel()
	wg.Wait()
}
