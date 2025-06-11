package domain

import "fmt"

// Secret represents a secret with name and value
type Secret struct {
	Name  string
	Value string
}

// EncryptedSecret represents an encrypted secret ready for storage
type EncryptedSecret struct {
	Name           string
	KeyID          string
	EncryptedValue string
}

// EncryptionService defines operations for secret encryption
type EncryptionService interface {
	Encrypt(plaintext []byte, publicKey *PublicKey) (string, error)
}

// NewSecret creates a new secret with validation
func NewSecret(name, value string) (*Secret, error) {
	if name == "" {
		return nil, fmt.Errorf("secret name cannot be empty")
	}
	if value == "" {
		return nil, fmt.Errorf("secret value cannot be empty")
	}
	return &Secret{
		Name:  name,
		Value: value,
	}, nil
}
