package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aereal/github-ops/internal/log"
)

func main() { os.Exit(run()) }

func run() int {
	log.Setup()
	ctx := context.Background()
	app := initializeApp()
	if err := app.Run(ctx, os.Args); err != nil {
		slog.ErrorContext(ctx, "Run failed", log.AttrError(err))
		return 1
	}
	return 0
}
