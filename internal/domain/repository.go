package domain

import (
	"strings"
)

type Repository struct {
	owner string
	name  string
}

func NewRepository(owner, name string) (*Repository, error) {
	if owner == "" {
		return nil, ErrEmptyRepositoryOwner
	}
	if name == "" {
		return nil, ErrEmptyRepositoryName
	}
	return &Repository{
		owner: owner,
		name:  name,
	}, nil
}

func ParseQualifiedRepository(qualified string) (*Repository, error) {
	owner, name, ok := strings.Cut(qualified, "/")
	if !ok {
		return nil, &MalformedQualifiedRepositoryError{Input: qualified}
	}
	return NewRepository(owner, name)
}

func (r *Repository) Owner() string {
	return r.owner
}

func (r *Repository) Name() string {
	return r.name
}

func (r *Repository) QualifiedName() string {
	return r.owner + "/" + r.name
}

func (r *Repository) String() string {
	return r.QualifiedName()
}

func (r *Repository) Equal(other *Repository) bool {
	if other == nil {
		return false
	}
	return r.owner == other.owner && r.name == other.name
}
