package tasks

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

type TaskStatus string

const (
	TaskQueued    TaskStatus = "queued"
	TaskRunning   TaskStatus = "running"
	TaskCompleted TaskStatus = "completed"
	TaskFailed    TaskStatus = "failed"
	TaskCancelled TaskStatus = "cancelled"
)

type TaskPriority int

const (
	TaskPriorityLow TaskPriority = iota
	TaskPriorityNormal
	TaskPriorityHigh
	TaskPriorityCritical
)

type Task struct {
	ID          string
	WorkspaceID string

	Type        string
	Name        string
	Description string

	Status   TaskStatus
	Priority TaskPriority

	Progress float64

	StartedAt  time.Time
	FinishedAt *time.Time

	CreatedBy string

	Payload map[string]any
	Result  map[string]any
	Error   string

	Metadata map[string]string

	// Advanced capabilities
	RetryCount   int
	MaxRetries   int
	Timeout      time.Duration
	Tags         []string
	ParentTaskID string
	ChildTaskIDs []string
}

// Interfaces to break import cycles and provide compile-time safety

type WorkspaceProvider interface {
	Overview(id string) (map[string]any, error)
	// Add other methods as needed by runners
}

type Indexer interface {
	Refresh(ctx context.Context, workspaceID string) error
	// Add other methods as needed by runners
}

type Store interface {
	// Define store methods as needed
}

type CommandsEngine interface {
	// Define command engine methods as needed
}

type CapabilitiesRegistry interface {
	// Define capability methods as needed
}

// Runtime is the shared context passed to all runners.
type Runtime struct {
	Logger zerolog.Logger

	Provider WorkspaceProvider
	Indexer  Indexer
	Store    Store
	EventBus *EventStream
	Commands CommandsEngine
	Registry CapabilitiesRegistry
}

// Runner is the pluggable execution model.
type Runner interface {
	Execute(ctx context.Context, task *Task, runtime *Runtime) error
}
