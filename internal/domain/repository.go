package domain

import (
	"context"
	"fmt"
	"strings"
)

// Repository represents a GitHub repository
type Repository struct {
	Owner string
	Name  string
}

// ParseQualifiedRepository parses a string in the format "owner/repo"
func ParseQualifiedRepository(s string) (*Repository, error) {
	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid repository format: expected 'owner/repo', got '%s'", s)
	}
	if parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid repository format: owner and repo name cannot be empty")
	}
	return &Repository{
		Owner: parts[0],
		Name:  parts[1],
	}, nil
}

// PublicKey represents a public key for encryption
type PublicKey struct {
	KeyID string
	Key   []byte
}

// RepositoryService defines operations for GitHub repositories
type RepositoryService interface {
	GetPublicKey(ctx context.Context, repo Repository) (*PublicKey, error)
	CreateOrUpdateSecret(ctx context.Context, repo Repository, secret EncryptedSecret) error
}
