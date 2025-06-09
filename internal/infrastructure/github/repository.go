package github

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/aereal/github-ops/internal/domain"
	"github.com/google/go-github/v72/github"
)

// RepositoryService implements domain.RepositoryService using GitHub API
type RepositoryService struct {
	client GHActionsService
}

// GHActionsService defines the interface for GitHub Actions API
type GHActionsService interface {
	GetRepoPublicKey(ctx context.Context, owner, repo string) (*github.PublicKey, *github.Response, error)
	CreateOrUpdateRepoSecret(ctx context.Context, owner, repo string, eSecret *github.EncryptedSecret) (*github.Response, error)
}

// NewRepositoryService creates a new GitHub repository service
func NewRepositoryService(client GHActionsService) *RepositoryService {
	return &RepositoryService{client: client}
}

// GetPublicKey retrieves the public key for a repository
func (s *RepositoryService) GetPublicKey(ctx context.Context, repo domain.Repository) (*domain.PublicKey, error) {
	pubKey, _, err := s.client.GetRepoPublicKey(ctx, repo.Owner, repo.Name)
	if err != nil {
		return nil, fmt.Errorf("GetRepoPublicKey: %w", err)
	}

	rawKey, err := base64.StdEncoding.DecodeString(pubKey.GetKey())
	if err != nil {
		return nil, fmt.Errorf("decode public key: %w", err)
	}

	return &domain.PublicKey{
		KeyID: pubKey.GetKeyID(),
		Key:   rawKey,
	}, nil
}

// CreateOrUpdateSecret creates or updates a secret in the repository
func (s *RepositoryService) CreateOrUpdateSecret(ctx context.Context, repo domain.Repository, secret domain.EncryptedSecret) error {
	ghSecret := &github.EncryptedSecret{
		Name:           secret.Name,
		KeyID:          secret.KeyID,
		EncryptedValue: secret.EncryptedValue,
	}

	_, err := s.client.CreateOrUpdateRepoSecret(ctx, repo.Owner, repo.Name, ghSecret)
	if err != nil {
		return fmt.Errorf("CreateOrUpdateRepoSecret: %w", err)
	}

	return nil
}
