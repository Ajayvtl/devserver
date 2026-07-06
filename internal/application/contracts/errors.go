package contracts

import "errors"

// Application-level Canonical Errors
var (
	ErrNotFound           = errors.New("application: resource not found")
	ErrValidationFailed   = errors.New("application: input validation failed")
	ErrPreconditionFailed = errors.New("application: precondition failed for operation")
	ErrConflict           = errors.New("application: state conflict")
	ErrUnauthorized       = errors.New("application: unauthorized access")
)
