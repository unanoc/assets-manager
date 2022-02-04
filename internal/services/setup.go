package services

import (
	log "github.com/sirupsen/logrus"

	"github.com/trustwallet/assets-manager/internal/config"
	"github.com/trustwallet/assets-manager/internal/queue"
	"github.com/trustwallet/go-libs/middleware"
)

func Setup() {
	// Init config.
	config.SetConfig()

	// Init logging.
	logLevel, err := log.ParseLevel(config.Default.LogLevel)
	if err != nil {
		log.WithError(err).Fatal("failed to parse log level")
	}

	log.SetLevel(logLevel)

	// Init sentry.
	err = middleware.SetupSentry(config.Default.Sentry.DSN,
		middleware.WithSampleRate(config.Default.Sentry.SampleRate))
	if err != nil {
		log.WithError(err).Error("failed to init Sentry")
	}

	// Init Rabbit MQ queues.
	err = queue.SetupQueues(config.Default.Rabbitmq.URL)
	if err != nil {
		log.WithError(err).Error("failed to init Rabbit MQ queues")
	}
}
