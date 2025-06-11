package encryption

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/aereal/github-ops/internal/domain"
	"golang.org/x/crypto/nacl/box"
)

// NaClService implements domain.EncryptionService using NaCl box encryption
type NaClService struct{}

// NewNaClService creates a new NaCl encryption service
func NewNaClService() *NaClService {
	return &NaClService{}
}

// Encrypt encrypts plaintext using the provided public key
func (s *NaClService) Encrypt(plaintext []byte, publicKey *domain.PublicKey) (string, error) {
	if len(publicKey.Key) != 32 {
		return "", fmt.Errorf("invalid public key length: expected 32 bytes, got %d", len(publicKey.Key))
	}

	// Convert to NaCl format
	var naclPubKey [32]byte
	if _, err := io.ReadFull(bytes.NewReader(publicKey.Key), naclPubKey[:]); err != nil {
		return "", fmt.Errorf("failed to read public key: %w", err)
	}

	// Encrypt using NaCl box
	var out []byte
	encrypted, err := box.SealAnonymous(out, plaintext, &naclPubKey, rand.Reader)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %w", err)
	}

	// Encode to base64
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// ProvideEncryptionService provides an EncryptionService instance
func ProvideEncryptionService() domain.EncryptionService {
	return NewNaClService()
}

var _ domain.EncryptionService = (*NaClService)(nil)
