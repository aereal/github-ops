package registersecret

type MissingTokenError struct{}

func (MissingTokenError) Error() string { return "missing GitHub token" }

var ErrMissingToken MissingTokenError
