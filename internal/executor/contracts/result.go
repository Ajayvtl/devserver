package contracts

import "time"

// ExecutionResult represents the canonical output of any executor.
type ExecutionResult struct {
	ExitCode    int
	Stdout      []byte
	Stderr      []byte
	Duration    time.Duration
	Metadata    map[string]interface{}
	Diagnostics []Diagnostic
}

// Diagnostic contains structured feedback or telemetry produced during execution.
type Diagnostic struct {
	Level   string // info, warn, error
	Message string
	Context map[string]interface{}
}
