package main

import (
	"context"

	"github.com/trustwallet/assets-manager/internal/services/consumer"
)

func main() {
	consumer.NewApp().Run(context.Background())
}
