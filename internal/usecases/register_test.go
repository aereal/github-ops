package usecases_test

import (
	"crypto/rand"
	"testing"

	"github.com/aereal/github-ops/internal/domain"
	"github.com/aereal/github-ops/internal/usecases"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/nacl/box"
)

func TestRegisterRepositorySecret_RegisterSecret(t *testing.T) {
	pubKey, err := getPublicKey()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		mockRepoService func(m *MockRepositoryService)
		mockEncryption  func(m *MockEncryptionService)
		request         domain.SecretRegistrationRequest
		name            string
		wantErr         bool
	}{
		{
			name: "ok",
			request: domain.SecretRegistrationRequest{
				Repository: domain.Repository{Owner: "aereal", Name: "myrepo"},
				Secret:     domain.Secret{Name: "MY_SECRET", Value: "blah blah"},
			},
			mockRepoService: func(m *MockRepositoryService) {
				m.EXPECT().GetPublicKey(gomock.Any(), gomock.Any()).Return(&domain.PublicKey{
					KeyID: "test-key-id",
					Key:   pubKey[:],
				}, nil)
				m.EXPECT().CreateOrUpdateSecret(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			mockEncryption: func(m *MockEncryptionService) {
				m.EXPECT().Encrypt(gomock.Any(), gomock.Any()).Return("encrypted-value", nil)
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockRepositoryService(ctrl)
			mockEncryption := NewMockEncryptionService(ctrl)

			tc.mockRepoService(mockRepo)
			tc.mockEncryption(mockEncryption)

			uc := usecases.NewRegisterRepositorySecret(mockRepo, mockEncryption)

			err := uc.RegisterSecret(t.Context(), tc.request)

			if tc.wantErr && err == nil {
				t.Errorf("expected error but got none")
			} else if !tc.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func getPublicKey() (*[32]byte, error) {
	pub, _, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return pub, nil
}
