package apperr

import (
	"fmt"
)

type AppError struct {
	Kind    Kind
	Code    Code
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError.
func New(kind Kind, code Code, msg string) error {
	return &AppError{Kind: kind, Code: code, Message: msg, Err: nil}
}

// Wraps an error with AppError.
func Wrap(err error, kind Kind, code Code, msg string) error {
	return &AppError{Kind: kind, Code: code, Message: msg, Err: err}
}

/* ==============================================================================
	 Specific Error Constructors
============================================================================== */

func NewValidationError(msg string, code ...Code) error {
	if len(code) > 0 {
		return New(Validation, code[0], msg)
	}

	return New(Validation, InvalidData, msg)
}

func NewRequestError(msg string, code ...Code) error {
	if len(code) > 0 {
		return New(Request, code[0], msg)
	}

	return New(Request, BadRequest, msg)
}

func NewUnauthorizedError(msg string, code ...Code) error {
	if len(code) > 0 {
		return New(Unauthorized, code[0], msg)
	}

	return New(Unauthorized, Unauthenticated, msg)
}

func NewForbiddenError(msg string, code ...Code) error {
	if len(code) > 0 {
		return New(Forbidden, code[0], msg)
	}

	return New(Forbidden, NotAllowed, msg)
}

func NewConflictError(msg string, code ...Code) error {
	if len(code) > 0 {
		return New(Conflict, code[0], msg)
	}

	return New(Conflict, Inconsistency, msg)
}

func NewInternalError(msg string, code ...Code) error {
	if len(code) > 0 {
		return New(Internal, code[0], msg)
	}

	return New(Internal, Unexpected, msg)
}

func NewExternalError(msg string, code ...Code) error {
	if len(code) > 0 {
		return New(External, code[0], msg)
	}

	return New(External, Unexpected, msg)
}
