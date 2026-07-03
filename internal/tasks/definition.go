// Package tasks provides a production-grade execution engine with
// priority queuing, worker pools, retry policies, and event-driven
// state transitions. It reuses core.Runner for rollback support.
package tasks

import (
	"time"
)

// Status represents the lifecycle state of a task.
type Status string

const (
	StatusPending    Status = "pending"
	StatusQueued     Status = "queued"
	StatusRunning    Status = "running"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
	StatusCancelled  Status = "cancelled"
	StatusRolledBack Status = "rolled_back"
)

// Priority determines queue ordering. Lower value = higher priority.
type Priority int

const (
	PriorityCritical Priority = 0
	PriorityHigh     Priority = 10
	PriorityNormal   Priority = 50
	PriorityLow      Priority = 100
)

// RetryPolicy controls how failed tasks are retried.
type RetryPolicy struct {
	MaxRetries int
	Backoff    time.Duration // Base backoff; doubles on each retry.
}

// Definition is the full specification of a task before execution.
type Definition struct {
	ID           string
	Name         string
	WorkspaceID  string
	Priority     Priority
	Retries      int
	Retry        RetryPolicy
	Timeout      time.Duration
	Rollback     string
	Dependencies []string
	Metadata     map[string]string
	// Fn is the work to execute. Return nil for success.
	Fn func(ctx *Context) error
	// RollbackFn is called if Fn fails or the task is rolled back.
	RollbackFn func(ctx *Context) error
}

// Record is the runtime state of a task at any point in its lifecycle.
type Record struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	WorkspaceID string            `json:"workspaceId"`
	Status      Status            `json:"status"`
	Progress    int               `json:"progress"`
	Detail      string            `json:"detail"`
	Priority    Priority          `json:"priority"`
	Attempt     int               `json:"attempt"`
	MaxRetries  int               `json:"maxRetries"`
	Error       string            `json:"error,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedAt   time.Time         `json:"createdAt"`
	StartedAt   *time.Time        `json:"startedAt,omitempty"`
	CompletedAt *time.Time        `json:"completedAt,omitempty"`
}

// Context is passed to task functions so they can report progress
// and check for cancellation without importing external packages.
type Context struct {
	record   *Record
	progress func(pct int, detail string)
	done     <-chan struct{}
}

// ReportProgress updates the task's progress percentage and detail message.
func (c *Context) ReportProgress(pct int, detail string) {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	c.record.Progress = pct
	c.record.Detail = detail
	if c.progress != nil {
		c.progress(pct, detail)
	}
}

// Done returns a channel that is closed when the task is cancelled.
func (c *Context) Done() <-chan struct{} {
	return c.done
}

// TaskID returns the current task's ID.
func (c *Context) TaskID() string {
	return c.record.ID
}
