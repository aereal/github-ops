//go:build wireinject

//go:generate go run -mod=mod github.com/google/wire/cmd/wire

package main

import (
	"github.com/aereal/github-ops/internal/cli/registersecret"
	"github.com/aereal/github-ops/internal/infrastructure/encryption"
	"github.com/aereal/github-ops/internal/infrastructure/github"
	"github.com/aereal/github-ops/internal/usecases"
	"github.com/google/wire"
)

// InfrastructureSet provides all infrastructure dependencies
var InfrastructureSet = wire.NewSet(
	github.ProvideGitHubClient,
	github.ProvideRepositoryService,
	encryption.ProvideEncryptionService,
)

// UseCaseSet provides all use case dependencies
var UseCaseSet = wire.NewSet(
	usecases.ProvideRegisterRepositorySecret,
)

// AppSet provides application layer dependencies
var AppSet = wire.NewSet(
	registersecret.ProvideApp,
)

// initializeApp initializes the application with all dependencies
func initializeApp() *registersecret.App {
	wire.Build(
		InfrastructureSet,
		UseCaseSet,
		AppSet,
	)
	return nil
}
