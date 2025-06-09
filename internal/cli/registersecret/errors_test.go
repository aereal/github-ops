package registersecret_test

import (
	"errors"
	"testing"

	"github.com/aereal/github-ops/internal/domain"
)

func TestMalformedQualifiedRepositoryError_Is(t *testing.T) {
	t.Run("same type, same value", func(t *testing.T) {
		this := &domain.MalformedQualifiedRepositoryError{Input: "repo1"}
		other := &domain.MalformedQualifiedRepositoryError{Input: "repo1"}
		if !errors.Is(this, other) {
			t.Error("errors.Is() expected to return true but got false")
		}
	})
	t.Run("same type, different value", func(t *testing.T) {
		this := &domain.MalformedQualifiedRepositoryError{Input: "repo1"}
		other := &domain.MalformedQualifiedRepositoryError{Input: "repo2"}
		if errors.Is(this, other) {
			t.Error("errors.Is() expected to return false but got true")
		}
	})
	t.Run("different type", func(t *testing.T) {
		this := &domain.MalformedQualifiedRepositoryError{Input: "repo1"}
		other := errors.New("oops")
		if errors.Is(this, other) {
			t.Error("errors.Is() expected to return false but got true")
		}
	})
}
