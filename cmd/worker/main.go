package main

import (
	"context"

	"github.com/trustwallet/assets-manager/internal/services/worker"
)

func main() {
	worker.NewApp().Run(context.Background())
}
