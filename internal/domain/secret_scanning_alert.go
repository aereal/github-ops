package domain

import (
	"context"
	"iter"
)

func DismissalIteratorFrom(alerts iter.Seq[*SecretScanningAlert], resolution SecretScanningAlertResolution, comment string) iter.Seq[*SecretScanningAlertDismissal] {
	return func(yield func(*SecretScanningAlertDismissal) bool) {
		for alert := range alerts {
			dismissal := &SecretScanningAlertDismissal{
				Alert:             alert,
				Resolution:        resolution,
				ResolutionComment: comment,
			}
			if !yield(dismissal) {
				return
			}
		}
	}
}

type SecretScanningAlertResolution int

const (
	SecretScanningAlertResolutionFalsePositive SecretScanningAlertResolution = iota + 1
	SecretScanningAlertResolutionWontFix
	SecretScanningAlertResolutionRevoked
	SecretScanningAlertResolutionUsedInTests
)

var resolution2str = map[SecretScanningAlertResolution]string{
	SecretScanningAlertResolutionFalsePositive: "false_positive",
	SecretScanningAlertResolutionWontFix:       "wont_fix",
	SecretScanningAlertResolutionRevoked:       "revoked",
	SecretScanningAlertResolutionUsedInTests:   "used_in_tests",
}

func (r SecretScanningAlertResolution) String() string {
	if msg, ok := resolution2str[r]; ok {
		return msg
	}
	return "INVALID"
}

func (r SecretScanningAlertResolution) Validate() error {
	if _, ok := resolution2str[r]; !ok {
		return ErrInvalidSecretScanningAlertResolution
	}
	return nil
}

var ErrInvalidSecretScanningAlertResolution InvalidSecretScanningAlertResolutionError

type InvalidSecretScanningAlertResolutionError struct{}

var _ error = InvalidSecretScanningAlertResolutionError{}

func (InvalidSecretScanningAlertResolutionError) Error() string {
	return "invalid secret scanning alert resolution"
}

type SecretScanningAlertDismissal struct {
	Alert             *SecretScanningAlert
	ResolutionComment string
	Resolution        SecretScanningAlertResolution
}

type SecretScanningAlertTargetRepository struct {
	Owner string
	Name  string
}

type SecretScanningAlert struct {
	Repository *SecretScanningAlertTargetRepository
	SecretType string
	Number     int64
}

type SecretScanningAlertIterator interface {
	Values(ctx context.Context) iter.Seq[*SecretScanningAlert]
	Err() error
}
