package internal

import (
	"context"
	"errors"
	"time"

	"github.com/Ajayvtl/devserver/internal/executor/contracts"
)

// NormalizeError converts raw OS and execution errors into canonical domain errors.
func NormalizeError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return contracts.ErrTimeout
	}
	if errors.Is(err, context.Canceled) {
		return contracts.ErrCancelled
	}
	if err.Error() == "executable file not found in %PATH%" || err.Error() == "executable file not found in $PATH" {
		return contracts.ErrBinaryMissing
	}
	return err
}

// Timeout execution context helper
func WithTimeout(ctx context.Context, duration time.Duration) (context.Context, context.CancelFunc) {
	if duration > 0 {
		return context.WithTimeout(ctx, duration)
	}
	return context.WithCancel(ctx)
}


