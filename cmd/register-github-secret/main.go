package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aereal/github-ops/internal/cli/registersecret"
	"github.com/aereal/github-ops/internal/infrastructure/encryption"
	"github.com/aereal/github-ops/internal/infrastructure/github"
	"github.com/aereal/github-ops/internal/log"
	"github.com/aereal/github-ops/internal/usecases"
	ghapi "github.com/google/go-github/v72/github"
)

func main() { os.Exit(run()) }

func run() int {
	log.Setup()
	ctx := context.Background()
	ghClient := ghapi.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN"))
	repoService := github.NewRepositoryService(ghClient.Actions)
	encryptionService := encryption.NewNaClService()
	uc := usecases.NewRegisterRepositorySecret(repoService, encryptionService)
	if err := registersecret.NewApp(uc).Run(ctx, os.Args); err != nil {
		slog.ErrorContext(ctx, "Run failed", log.AttrError(err))
		return 1
	}
	return 0
}
