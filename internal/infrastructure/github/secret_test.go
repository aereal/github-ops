package github

import (
	"context"
	"errors"
	"testing"

	"github.com/aereal/github-ops/internal/domain"
	"github.com/aereal/github-ops/internal/infrastructure/github/mocks"
	"github.com/google/go-github/v72/github"
	"go.uber.org/mock/gomock"
)

func TestNewSecretScanAlertService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockGHSecretScanningService(ctrl)
	service := NewSecretScanAlertService(client)

	if service.client != client {
		t.Errorf("Expected client to be set correctly")
	}
}

func TestSecretScanAlertService_ListOpenSecretScanningAlertsForOrg(t *testing.T) {
	tests := []struct {
		name        string
		org         string
		secretTypes []string
		mockSetup   func(*mocks.MockGHSecretScanningService)
		expected    []*domain.SecretScanningAlert
		expectError bool
	}{
		{
			name:        "successful single page response",
			org:         "testorg",
			secretTypes: []string{"github_token"},
			mockSetup: func(m *mocks.MockGHSecretScanningService) {
				alerts := []*github.SecretScanningAlert{
					{
						Number:     github.Ptr(1),
						SecretType: github.Ptr("github_token"),
					},
					{
						Number:     github.Ptr(2),
						SecretType: github.Ptr("github_token"),
					},
				}
				resp := &github.Response{After: ""}
				m.EXPECT().ListAlertsForOrg(
					gomock.Any(),
					"testorg",
					gomock.Any(),
				).Return(alerts, resp, nil)
			},
			expected: []*domain.SecretScanningAlert{
				{Number: 1, SecretType: "github_token"},
				{Number: 2, SecretType: "github_token"},
			},
			expectError: false,
		},
		{
			name:        "successful multi-page response",
			org:         "testorg",
			secretTypes: []string{"aws_access_key"},
			mockSetup: func(m *mocks.MockGHSecretScanningService) {
				// First page
				alerts1 := []*github.SecretScanningAlert{
					{
						Number:     github.Ptr(1),
						SecretType: github.Ptr("aws_access_key"),
					},
				}
				resp1 := &github.Response{After: "page2"}

				// Second page
				alerts2 := []*github.SecretScanningAlert{
					{
						Number:     github.Ptr(2),
						SecretType: github.Ptr("aws_access_key"),
					},
				}
				resp2 := &github.Response{After: ""}

				gomock.InOrder(
					m.EXPECT().ListAlertsForOrg(
						gomock.Any(),
						"testorg",
						gomock.Any(),
					).Return(alerts1, resp1, nil),
					m.EXPECT().ListAlertsForOrg(
						gomock.Any(),
						"testorg",
						gomock.Any(),
					).Return(alerts2, resp2, nil),
				)
			},
			expected: []*domain.SecretScanningAlert{
				{Number: 1, SecretType: "aws_access_key"},
				{Number: 2, SecretType: "aws_access_key"},
			},
			expectError: false,
		},
		{
			name:        "empty response",
			org:         "testorg",
			secretTypes: []string{"slack_token"},
			mockSetup: func(m *mocks.MockGHSecretScanningService) {
				alerts := []*github.SecretScanningAlert{}
				resp := &github.Response{After: ""}
				m.EXPECT().ListAlertsForOrg(
					gomock.Any(),
					"testorg",
					gomock.Any(),
				).Return(alerts, resp, nil)
			},
			expected:    []*domain.SecretScanningAlert{},
			expectError: false,
		},
		{
			name:        "API error",
			org:         "testorg",
			secretTypes: []string{"github_token"},
			mockSetup: func(m *mocks.MockGHSecretScanningService) {
				m.EXPECT().ListAlertsForOrg(
					gomock.Any(),
					"testorg",
					gomock.Any(),
				).Return(nil, nil, errors.New("API error"))
			},
			expected:    []*domain.SecretScanningAlert{},
			expectError: true,
		},
		{
			name:        "multiple secret types",
			org:         "testorg",
			secretTypes: []string{"github_token", "aws_access_key"},
			mockSetup: func(m *mocks.MockGHSecretScanningService) {
				alerts := []*github.SecretScanningAlert{
					{
						Number:     github.Ptr(1),
						SecretType: github.Ptr("github_token"),
					},
					{
						Number:     github.Ptr(2),
						SecretType: github.Ptr("aws_access_key"),
					},
				}
				resp := &github.Response{After: ""}
				m.EXPECT().ListAlertsForOrg(
					gomock.Any(),
					"testorg",
					gomock.Any(),
				).Return(alerts, resp, nil)
			},
			expected: []*domain.SecretScanningAlert{
				{Number: 1, SecretType: "github_token"},
				{Number: 2, SecretType: "aws_access_key"},
			},
			expectError: false,
		},
		{
			name:        "nil values in GitHub response",
			org:         "testorg",
			secretTypes: []string{"github_token"},
			mockSetup: func(m *mocks.MockGHSecretScanningService) {
				alerts := []*github.SecretScanningAlert{
					{
						Number:     nil,
						SecretType: nil,
					},
					{
						Number:     github.Ptr(1),
						SecretType: github.Ptr("github_token"),
					},
				}
				resp := &github.Response{After: ""}
				m.EXPECT().ListAlertsForOrg(
					gomock.Any(),
					"testorg",
					gomock.Any(),
				).Return(alerts, resp, nil)
			},
			expected: []*domain.SecretScanningAlert{
				{Number: 0, SecretType: ""},
				{Number: 1, SecretType: "github_token"},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockGHSecretScanningService(ctrl)
			tt.mockSetup(mockClient)

			service := &SecretScanAlertService{client: mockClient}
			iterator := service.ListOpenSecretScanningAlertsForOrg(t.Context(), tt.org, tt.secretTypes)

			var results []*domain.SecretScanningAlert
			for alert := range iterator.Values(t.Context()) {
				results = append(results, alert)
			}

			if tt.expectError {
				if iterator.Err() == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if iterator.Err() != nil {
					t.Errorf("Unexpected error: %v", iterator.Err())
				}
			}

			if len(results) != len(tt.expected) {
				t.Errorf("Expected %d results, got %d", len(tt.expected), len(results))
				return
			}

			for i, result := range results {
				if result.Number != tt.expected[i].Number {
					t.Errorf("Expected Number %d, got %d", tt.expected[i].Number, result.Number)
				}
				if result.SecretType != tt.expected[i].SecretType {
					t.Errorf("Expected SecretType %s, got %s", tt.expected[i].SecretType, result.SecretType)
				}
			}
		})
	}
}

func TestAlertIterator_Err(t *testing.T) {
	tests := []struct {
		setupError  error
		expectedErr error
		name        string
	}{
		{
			name:        "no error",
			setupError:  nil,
			expectedErr: nil,
		},
		{
			name:        "with error",
			setupError:  errors.New("test error"),
			expectedErr: errors.New("test error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iterator := &alertIterator{lastErr: tt.setupError}

			err := iterator.Err()

			if tt.expectedErr == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if tt.expectedErr != nil && err == nil {
				t.Errorf("Expected error %v, got nil", tt.expectedErr)
			}
			if tt.expectedErr != nil && err != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}

func TestAlertIterator_Values_EarlyBreak(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockGHSecretScanningService(ctrl)

	alerts := []*github.SecretScanningAlert{
		{
			Number:     github.Ptr(1),
			SecretType: github.Ptr("github_token"),
		},
		{
			Number:     github.Ptr(2),
			SecretType: github.Ptr("github_token"),
		},
		{
			Number:     github.Ptr(3),
			SecretType: github.Ptr("github_token"),
		},
	}
	resp := &github.Response{After: ""}

	mockClient.EXPECT().ListAlertsForOrg(
		gomock.Any(),
		"testorg",
		gomock.Any(),
	).Return(alerts, resp, nil)

	service := &SecretScanAlertService{client: mockClient}
	iterator := service.ListOpenSecretScanningAlertsForOrg(t.Context(), "testorg", []string{"github_token"})

	var results []*domain.SecretScanningAlert
	count := 0
	for alert := range iterator.Values(t.Context()) {
		results = append(results, alert)
		count++
		if count >= 2 {
			break
		}
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results due to early break, got %d", len(results))
	}

	if iterator.Err() != nil {
		t.Errorf("Unexpected error: %v", iterator.Err())
	}
}

func TestSecretScanAlertService_OptionsConfiguration(t *testing.T) {
	tests := []struct {
		expectedOpts func(t *testing.T, opts *github.SecretScanningAlertListOptions)
		name         string
		org          string
		secretTypes  []string
	}{
		{
			name:        "single secret type",
			org:         "testorg",
			secretTypes: []string{"github_token"},
			expectedOpts: func(t *testing.T, opts *github.SecretScanningAlertListOptions) {
				t.Helper()
				if opts.State != "open" {
					t.Errorf("Expected state to be 'open', got '%s'", opts.State)
				}
				if opts.SecretType != "github_token" {
					t.Errorf("Expected SecretType to be 'github_token', got '%s'", opts.SecretType)
				}
			},
		},
		{
			name:        "multiple secret types",
			org:         "testorg",
			secretTypes: []string{"github_token", "aws_access_key", "slack_token"},
			expectedOpts: func(t *testing.T, opts *github.SecretScanningAlertListOptions) {
				t.Helper()
				if opts.State != "open" {
					t.Errorf("Expected state to be 'open', got '%s'", opts.State)
				}
				expected := "github_token,aws_access_key,slack_token"
				if opts.SecretType != expected {
					t.Errorf("Expected SecretType to be '%s', got '%s'", expected, opts.SecretType)
				}
			},
		},
		{
			name:        "empty secret types",
			org:         "testorg",
			secretTypes: []string{},
			expectedOpts: func(t *testing.T, opts *github.SecretScanningAlertListOptions) {
				t.Helper()
				if opts.State != "open" {
					t.Errorf("Expected state to be 'open', got '%s'", opts.State)
				}
				if opts.SecretType != "" {
					t.Errorf("Expected SecretType to be empty, got '%s'", opts.SecretType)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClient := mocks.NewMockGHSecretScanningService(ctrl)

			// Capture the options passed to the mock
			var capturedOpts *github.SecretScanningAlertListOptions
			mockClient.EXPECT().ListAlertsForOrg(
				gomock.Any(),
				tt.org,
				gomock.Any(),
			).DoAndReturn(func(ctx context.Context, org string, opts *github.SecretScanningAlertListOptions) ([]*github.SecretScanningAlert, *github.Response, error) {
				capturedOpts = opts
				return []*github.SecretScanningAlert{}, &github.Response{After: ""}, nil
			})

			service := &SecretScanAlertService{client: mockClient}
			iterator := service.ListOpenSecretScanningAlertsForOrg(t.Context(), tt.org, tt.secretTypes)

			// Consume the iterator to trigger the API call
			for range iterator.Values(t.Context()) {
				break
			}

			tt.expectedOpts(t, capturedOpts)
		})
	}
}

func TestAlertIterator_InterfaceCompliance(t *testing.T) {
	// Ensure alertIterator implements domain.SecretScanningAlertIterator
	var _ domain.SecretScanningAlertIterator = (*alertIterator)(nil)
}

func TestSecretScanAlertService_EmptySecretTypes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockGHSecretScanningService(ctrl)

	mockClient.EXPECT().ListAlertsForOrg(
		gomock.Any(),
		"testorg",
		gomock.Any(),
	).Return([]*github.SecretScanningAlert{}, &github.Response{After: ""}, nil)

	service := &SecretScanAlertService{client: mockClient}
	iterator := service.ListOpenSecretScanningAlertsForOrg(t.Context(), "testorg", []string{})

	var results []*domain.SecretScanningAlert
	for alert := range iterator.Values(t.Context()) {
		results = append(results, alert)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results for empty secret types, got %d", len(results))
	}

	if iterator.Err() != nil {
		t.Errorf("Unexpected error: %v", iterator.Err())
	}
}

func TestAlertIterator_Values_YieldReturnsFalse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockGHSecretScanningService(ctrl)

	alerts := []*github.SecretScanningAlert{
		{
			Number:     github.Ptr(1),
			SecretType: github.Ptr("github_token"),
		},
		{
			Number:     github.Ptr(2),
			SecretType: github.Ptr("github_token"),
		},
		{
			Number:     github.Ptr(3),
			SecretType: github.Ptr("github_token"),
		},
	}
	resp := &github.Response{After: ""}

	mockClient.EXPECT().ListAlertsForOrg(
		gomock.Any(),
		"testorg",
		gomock.Any(),
	).Return(alerts, resp, nil)

	service := &SecretScanAlertService{client: mockClient}
	iterator := service.ListOpenSecretScanningAlertsForOrg(t.Context(), "testorg", []string{"github_token"})

	// Manually call the iterator with a yield function that returns false after the first call
	var results []*domain.SecretScanningAlert
	count := 0
	seq := iterator.Values(t.Context())
	seq(func(alert *domain.SecretScanningAlert) bool {
		results = append(results, alert)
		count++
		return count < 1 // Return false after first element
	})

	if len(results) != 1 {
		t.Errorf("Expected 1 result due to yield returning false, got %d", len(results))
	}

	if iterator.Err() != nil {
		t.Errorf("Unexpected error: %v", iterator.Err())
	}
}

func TestAlertIterator_Values_MultiPageWithYieldFalse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockGHSecretScanningService(ctrl)

	// First page - should be called
	alerts1 := []*github.SecretScanningAlert{
		{
			Number:     github.Ptr(1),
			SecretType: github.Ptr("github_token"),
		},
		{
			Number:     github.Ptr(2),
			SecretType: github.Ptr("github_token"),
		},
	}
	resp1 := &github.Response{After: "page2"}

	// Second page - should NOT be called because yield returns false
	mockClient.EXPECT().ListAlertsForOrg(
		gomock.Any(),
		"testorg",
		gomock.Any(),
	).Return(alerts1, resp1, nil)

	service := &SecretScanAlertService{client: mockClient}
	iterator := service.ListOpenSecretScanningAlertsForOrg(t.Context(), "testorg", []string{"github_token"})

	// Manually call the iterator with a yield function that returns false after the first call
	var results []*domain.SecretScanningAlert
	count := 0
	seq := iterator.Values(t.Context())
	seq(func(alert *domain.SecretScanningAlert) bool {
		results = append(results, alert)
		count++
		return count < 1 // Return false after first element to trigger early return
	})

	if len(results) != 1 {
		t.Errorf("Expected 1 result due to yield returning false, got %d", len(results))
	}

	if iterator.Err() != nil {
		t.Errorf("Unexpected error: %v", iterator.Err())
	}
}
