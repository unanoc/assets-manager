package main

import (
	"context"

	"github.com/trustwallet/assets-manager/internal/services/api"
)

func main() {
	api.NewApp().Run(context.Background())
}
