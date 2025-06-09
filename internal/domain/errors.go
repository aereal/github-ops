package domain

import (
	"errors"
	"fmt"
)

var (
	ErrEmptyRepositoryOwner    = errors.New("repository owner cannot be empty")
	ErrEmptyRepositoryName     = errors.New("repository name cannot be empty")
	ErrEmptySecretName         = errors.New("secret name cannot be empty")
	ErrEmptySecretValue        = errors.New("secret value cannot be empty")
	ErrNoRepositoriesSpecified = errors.New("no repositories specified")
	ErrSecretRequired          = errors.New("secret is required")
)

type MalformedQualifiedRepositoryError struct {
	Input string
}

func (e *MalformedQualifiedRepositoryError) Error() string {
	return fmt.Sprintf("malformed qualified repository name: %q", e.Input)
}

func (e *MalformedQualifiedRepositoryError) Is(err error) bool {
	var target *MalformedQualifiedRepositoryError
	if !errors.As(err, &target) {
		return false
	}
	return e.Input == target.Input
}
