package queue

import (
	"fmt"

	"github.com/trustwallet/go-libs/mq"
)

const (
	QueueAssetManagerProcessGithubEvent mq.QueueName = "asset_manager_github_events_process"
)

func SetupQueues(rabbitmqURL string) error {
	mqClient, err := mq.Connect(rabbitmqURL)
	if err != nil {
		return fmt.Errorf("failed to connect to mq: %w", err)
	}

	queues := []mq.Queue{
		mqClient.InitQueue(QueueAssetManagerProcessGithubEvent),
	}

	for _, queue := range queues {
		if err := queue.Declare(); err != nil {
			return fmt.Errorf("failed to declare queue (%s): %w", queue.Name(), err)
		}
	}

	return nil
}
