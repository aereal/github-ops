package domain

import (
	"context"
)

type SecretRegistrationRequest struct {
	secret       *Secret
	repositories []*Repository
}

func NewSecretRegistrationRequest(repositories []*Repository, secret *Secret) (*SecretRegistrationRequest, error) {
	if len(repositories) == 0 {
		return nil, ErrNoRepositoriesSpecified
	}
	if secret == nil {
		return nil, ErrSecretRequired
	}
	return &SecretRegistrationRequest{
		repositories: repositories,
		secret:       secret,
	}, nil
}

func (r *SecretRegistrationRequest) Repositories() []*Repository {
	return r.repositories
}

func (r *SecretRegistrationRequest) Secret() *Secret {
	return r.secret
}

type SecretRegistrationService interface {
	RegisterSecret(ctx context.Context, request *SecretRegistrationRequest) error
}
