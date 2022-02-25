package metrics

import (
	"fmt"
	"time"

	"github.com/trustwallet/go-libs/metrics"
	"github.com/trustwallet/go-libs/worker"
)

func InitMetricsPusher(pushgatewayURL string, pushInterval time.Duration) (worker.Worker, error) {
	pusher := metrics.NewPusher(pushgatewayURL, "assets_manager_worker")

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
