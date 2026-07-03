package tasks

import "time"

const defaultRetryBackoff = 500 * time.Millisecond

// maxRetries resolves the retry budget from the legacy RetryPolicy
// or the direct Retries field.
func maxRetries(def *Definition) int {
	if def == nil {
		return 0
	}
	if def.Retry.MaxRetries > 0 {
		return def.Retry.MaxRetries
	}
	return def.Retries
}

// retryBackoff returns the backoff for a given attempt number.
func retryBackoff(def *Definition, attempt int) time.Duration {
	if def == nil {
		return defaultRetryBackoff
	}
	base := def.Retry.Backoff
	if base <= 0 {
		base = defaultRetryBackoff
	}
	if attempt <= 1 {
		return base
	}
	return base * time.Duration(1<<(attempt-1))
}
