package metrics

import (
	"fmt"
	"time"

	"github.com/trustwallet/go-libs/metrics"
	"github.com/trustwallet/go-libs/worker"
	metricspusher "github.com/trustwallet/go-libs/worker/metrics"
)

func InitMetricsPusher(pushgatewayURL, pushgatewayKey string, pushInterval time.Duration) (worker.Worker, error) {
	client := metrics.NewMetricsPusherClient(pushgatewayURL, pushgatewayKey, nil)
	pusher := metrics.NewPusherWithCustomClient(pushgatewayURL, "assets_manager_consumer", client)

	err := pusher.Push()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Pushgateway, metrics won't be pushed: %w", err)
	}

	metricsPusher := metricspusher.NewMetricsPusherWorker(worker.DefaultWorkerOptions(pushInterval), pusher)

	return metricsPusher, nil
}
