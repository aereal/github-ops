//go:generate go tool mockgen -destination ./usecase_mock_test.go -package registersecret_test -typed -write_command_comment=false github.com/aereal/github-ops/internal/domain SecretRegistrationService

package registersecret

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/aereal/github-ops/internal/domain"
	set "github.com/hashicorp/go-set/v3"
	"golang.org/x/sync/errgroup"
)

type SecretRegistrationUsecase = domain.SecretRegistrationService

func NewApp(uc SecretRegistrationUsecase) *App {
	return &App{uc: uc}
}

type App struct {
	uc SecretRegistrationUsecase
}

func (a *App) Run(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet(filepath.Base(args[0]), flag.ContinueOnError)
	var (
		secretName  string
		secretValue string
		repos       = set.New[string](0)
	)
	fs.Func("repos", "repository name list", func(s string) error {
		_, err := domain.ParseQualifiedRepository(s)
		if err != nil {
			return err
		}
		_ = repos.Insert(s)
		return nil
	})
	fs.StringVar(&secretName, "secret-name", "", "secret name")
	fs.StringVar(&secretValue, "secret-value", "", "secret value")
	err := fs.Parse(args[1:])
	switch {
	case errors.Is(err, flag.ErrHelp):
		return nil
	case err != nil:
		return err
	}
	secret, err := domain.NewSecret(secretName, secretValue)
	if err != nil {
		return err
	}
	eg, ctx := errgroup.WithContext(ctx)
	for repoStr := range repos.Items() {
		repo, err := domain.ParseQualifiedRepository(repoStr)
		if err != nil {
			return err
		}
		eg.Go(func() error {
			request := domain.NewSecretRegistrationRequest(*repo, *secret)
			return a.uc.RegisterSecret(ctx, request)
		})
	}
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("usecases.NewRegisterRepositorySecret.Do: %w", err)
	}
	return nil
}

// ProvideApp provides an App instance
func ProvideApp(uc SecretRegistrationUsecase) *App {
	return NewApp(uc)
}
