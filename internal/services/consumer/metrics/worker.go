package metrics

import (
	"fmt"
	"time"

	"github.com/trustwallet/go-libs/client"
	"github.com/trustwallet/go-libs/metrics"
	"github.com/trustwallet/go-libs/worker"
)

func InitMetricsPusher(pushgatewayURL, pushgatewayKey string, pushInterval time.Duration) (worker.Worker, error) {
	client := client.InitClient(pushgatewayURL, nil)
	client.AddHeader("X-API-Key", pushgatewayKey)

	pusher := metrics.NewPusherWithCustomClient(pushgatewayURL, "assets_manager_consumer", client.HttpClient)

	// Check connection to pusher.
	err := pusher.Push()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Pushgateway, metrics won't be pushed: %w", err)
	}

	metricsPusher := worker.InitWorker(
		"metrics_pusher",
		worker.DefaultWorkerOptions(pushInterval),
		pusher.Push,
	)

	return metricsPusher, nil
}
