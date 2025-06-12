package domain

import (
	"context"
	"iter"
)

type SecretScanningAlert struct {
	SecretType string
	Number     int
}

type SecretScanningAlertIterator interface {
	Values(ctx context.Context) iter.Seq[*SecretScanningAlert]
	Err() error
}
