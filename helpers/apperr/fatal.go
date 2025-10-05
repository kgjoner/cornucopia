package apperr

import "errors"

type fatalError struct {
	Err error
}

func (e *fatalError) Error() string {
	return e.Err.Error()
}

func (e *fatalError) Unwrap() error {
	return e.Err
}

// Fatal wraps an error to mark it as fatal.
func Fatal(err error) error {
	return &fatalError{Err: err}
}

// IsFatal checks if any error in the chain is marked as fatal.
func IsFatal(err error) bool {
	var fatalErr *fatalError
	return errors.As(err, &fatalErr)
}
