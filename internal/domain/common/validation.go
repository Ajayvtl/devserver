package common

import (
	"errors"
	"strings"
)

var (
	ErrInvalidName      = errors.New("invalid name: must not be empty and cannot contain special characters")
	ErrInvalidPath      = errors.New("invalid path: must not be empty")
	ErrInvalidEndpoint  = errors.New("invalid endpoint: must be a valid URL or host:port")
	ErrInvalidVersion   = errors.New("invalid version: must not be empty")
	ErrInvalidReference = errors.New("invalid reference: must not be empty")
)

// ValidateName ensures names conform to canonical rules.
func ValidateName(name string) error {
	trimmed := strings.TrimSpace(name)
	if len(trimmed) == 0 {
		return ErrInvalidName
	}
	return nil
}

func ValidatePath(p string) error {
	if len(strings.TrimSpace(p)) == 0 {
		return ErrInvalidPath
	}
	return nil
}

func ValidateEndpoint(e string) error {
	if len(strings.TrimSpace(e)) == 0 {
		return ErrInvalidEndpoint
	}
	return nil
}

func ValidateVersion(v string) error {
	if len(strings.TrimSpace(v)) == 0 {
		return ErrInvalidVersion
	}
	return nil
}

func ValidateReference(r string) error {
	if len(strings.TrimSpace(r)) == 0 {
		return ErrInvalidReference
	}
	return nil
}
