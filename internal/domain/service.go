package domain

import "context"

// SecretRegistrationRequest represents a request to register a secret
type SecretRegistrationRequest struct {
	Repository Repository
	Secret     Secret
}

// SecretRegistrationService defines the business logic for secret registration
type SecretRegistrationService interface {
	RegisterSecret(ctx context.Context, request SecretRegistrationRequest) error
}

// SecretScanningAlertService defines the business logic for secret scanning alerts
type SecretScanningAlertService interface {
	ListOpenSecretScanningAlertsForOrg(ctx context.Context, org string, secretTypes []string) SecretScanningAlertIterator
}

// NewSecretRegistrationRequest creates a new secret registration request
func NewSecretRegistrationRequest(repo Repository, secret Secret) SecretRegistrationRequest {
	return SecretRegistrationRequest{
		Repository: repo,
		Secret:     secret,
	}
}
