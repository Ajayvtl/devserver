package contracts

import "errors"

// Canonical executor errors prevent raw OS/Network errors from crossing boundaries.
var (
	ErrNotSupported     = errors.New("executor: action not supported by capabilities")
	ErrPermissionDenied = errors.New("executor: permission denied")
	ErrTimeout          = errors.New("executor: operation timed out")
	ErrConnectionFailed = errors.New("executor: connection to target failed")
	ErrBinaryMissing    = errors.New("executor: required binary missing")
	ErrCancelled        = errors.New("executor: operation cancelled")
	ErrValidationFailed = errors.New("executor: action validation failed")
)
