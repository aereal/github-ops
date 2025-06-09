package domain

type Secret struct {
	name  string
	value string
}

func NewSecret(name, value string) (*Secret, error) {
	if name == "" {
		return nil, ErrEmptySecretName
	}
	if value == "" {
		return nil, ErrEmptySecretValue
	}
	return &Secret{
		name:  name,
		value: value,
	}, nil
}

func (s *Secret) Name() string {
	return s.name
}

func (s *Secret) Value() string {
	return s.value
}
