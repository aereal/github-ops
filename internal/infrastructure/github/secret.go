package github

import (
	"context"
	"iter"
	"strings"

	"github.com/aereal/github-ops/internal/domain"
	"github.com/google/go-github/v72/github"
)

type GHSecretScanningService interface {
	ListAlertsForOrg(ctx context.Context, org string, opts *github.SecretScanningAlertListOptions) ([]*github.SecretScanningAlert, *github.Response, error)
}

type SecretScanAlertService struct {
	client GHSecretScanningService
}

var _ domain.SecretScanningAlertService = (*SecretScanAlertService)(nil)

// NewSecretScanAlertService creates a new SecretScanAlertService instance.
func NewSecretScanAlertService(client GHSecretScanningService) *SecretScanAlertService {
	return &SecretScanAlertService{
		client: client,
	}
}

func (s *SecretScanAlertService) ListOpenSecretScanningAlertsForOrg(ctx context.Context, org string, secretTypes []string) domain.SecretScanningAlertIterator {
	i := &alertIterator{
		svc: s.client,
		opts: &github.SecretScanningAlertListOptions{
			State:      "open",
			SecretType: strings.Join(secretTypes, ","),
		},
		org:     org,
		lastErr: nil,
	}
	return i
}

type alertIterator struct {
	svc     GHSecretScanningService
	lastErr error
	opts    *github.SecretScanningAlertListOptions
	org     string
}

var _ domain.SecretScanningAlertIterator = (*alertIterator)(nil)

func (i *alertIterator) Err() error { return i.lastErr }

func (i *alertIterator) Values(ctx context.Context) iter.Seq[*domain.SecretScanningAlert] {
	return func(yield func(*domain.SecretScanningAlert) bool) {
		for {
			alerts, resp, err := i.svc.ListAlertsForOrg(ctx, i.org, i.opts)
			if err != nil {
				i.lastErr = err
				return
			}
			for _, a := range alerts {
				alert := &domain.SecretScanningAlert{
					Number:     a.GetNumber(),
					SecretType: a.GetSecretType(),
				}
				if !yield(alert) {
					return
				}
			}
			i.opts.After = resp.After
			if i.opts.After == "" {
				break
			}
		}
	}
}
