package main

import (
	"context"

	app "github.com/trustwallet/assets-manager/internal/merge-fee-bot"
)

func main() {
	app.NewApp().Run(context.Background())
}
