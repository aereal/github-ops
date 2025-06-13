package domain

import (
	"context"
	"iter"
)

type SecretScanningAlertTargetRepository struct {
	Owner string
	Name  string
}

type SecretScanningAlert struct {
	Repository *SecretScanningAlertTargetRepository
	SecretType string
	Number     int
}

type SecretScanningAlertIterator interface {
	Values(ctx context.Context) iter.Seq[*SecretScanningAlert]
	Err() error
}
