package github

import (
	"context"
	"iter"
	"strings"

	"github.com/aereal/github-ops/internal/domain"
	"github.com/google/go-github/v72/github"
	"golang.org/x/sync/errgroup"
)

type GHSecretScanningService interface {
	ListAlertsForOrg(ctx context.Context, org string, opts *github.SecretScanningAlertListOptions) ([]*github.SecretScanningAlert, *github.Response, error)
	UpdateAlert(ctx context.Context, owner, repo string, number int64, opts *github.SecretScanningAlertUpdateOptions) (*github.SecretScanningAlert, *github.Response, error)
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

func (s *SecretScanAlertService) DismissSecretScanningAlerts(ctx context.Context, dismissals iter.Seq[*domain.SecretScanningAlertDismissal]) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(2) // TODO
	for dismissal := range dismissals {
		eg.Go(func() error {
			return s.dismissSecretScanningAlert(ctx, dismissal)
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (s *SecretScanAlertService) dismissSecretScanningAlert(ctx context.Context, dismissal *domain.SecretScanningAlertDismissal) error {
	opts := &github.SecretScanningAlertUpdateOptions{
		State:             "resolved",
		Resolution:        github.Ptr(dismissal.Resolution.String()),
		ResolutionComment: github.Ptr(dismissal.ResolutionComment),
	}
	alert := dismissal.Alert
	repo := alert.Repository
	if _, _, err := s.client.UpdateAlert(ctx, repo.Owner, repo.Name, alert.Number, opts); err != nil {
		return err
	}
	return nil
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
					Number:     int64(a.GetNumber()),
					SecretType: a.GetSecretType(),
					Repository: &domain.SecretScanningAlertTargetRepository{
						Owner: a.GetRepository().GetOwner().GetLogin(),
						Name:  a.GetRepository().GetName(),
					},
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
