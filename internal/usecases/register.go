//go:generate go tool mockgen -destination ./mock_test.go -package usecases_test -typed -write_command_comment=false github.com/aereal/github-ops/internal/domain RepositoryService,EncryptionService

package usecases

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aereal/github-ops/internal/domain"
)

func NewRegisterRepositorySecret(repoService domain.RepositoryService, encryptionService domain.EncryptionService) *RegisterRepositorySecret {
	return &RegisterRepositorySecret{
		repoService:       repoService,
		encryptionService: encryptionService,
	}
}

type RegisterRepositorySecret struct {
	repoService       domain.RepositoryService
	encryptionService domain.EncryptionService
}

func (u *RegisterRepositorySecret) RegisterSecret(ctx context.Context, request domain.SecretRegistrationRequest) error {
	// Get public key for the repository
	pubKey, err := u.repoService.GetPublicKey(ctx, request.Repository)
	if err != nil {
		return fmt.Errorf("failed to get public key: %w", err)
	}

	// Encrypt the secret value
	encryptedValue, err := u.encryptionService.Encrypt([]byte(request.Secret.Value), pubKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt secret: %w", err)
	}

	// Create encrypted secret
	encryptedSecret := domain.EncryptedSecret{
		Name:           request.Secret.Name,
		KeyID:          pubKey.KeyID,
		EncryptedValue: encryptedValue,
	}

	slog.InfoContext(ctx, "set repository secret",
		slog.String("repo.owner", request.Repository.Owner),
		slog.String("repo.name", request.Repository.Name),
		slog.String("secret.name", request.Secret.Name),
	)

	// Store the encrypted secret
	if err := u.repoService.CreateOrUpdateSecret(ctx, request.Repository, encryptedSecret); err != nil {
		return fmt.Errorf("failed to create or update secret: %w", err)
	}

	return nil
}

// ProvideRegisterRepositorySecret provides a RegisterRepositorySecret instance
func ProvideRegisterRepositorySecret(repoService domain.RepositoryService, encryptionService domain.EncryptionService) domain.SecretRegistrationService {
	return NewRegisterRepositorySecret(repoService, encryptionService)
}

var _ domain.SecretRegistrationService = (*RegisterRepositorySecret)(nil)
