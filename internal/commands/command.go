package commands

import (
	"time"

	"github.com/Ajayvtl/devserver/internal/capabilities"
)

// Status tracks the lifecycle of a command.
type Status string

const (
	StatusQueued     Status = "queued"
	StatusRunning    Status = "running"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
	StatusCancelled  Status = "cancelled"
	StatusRolledBack Status = "rolled_back"
)

// Command is the user-facing intent that enters the system.
type Command struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	WorkspaceID string                  `json:"workspaceId,omitempty"`
	Capability  capabilities.Capability `json:"capability"`
	Provider    string                  `json:"provider,omitempty"`
	Target      string                  `json:"target,omitempty"`
	Parameters  map[string]any          `json:"parameters,omitempty"`
	Metadata    map[string]string       `json:"metadata,omitempty"`
	Retries     int                     `json:"retries,omitempty"`
	Timeout     time.Duration           `json:"timeout,omitempty"`
}

// Record captures the command lifecycle for UI/API inspection.
type Record struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	WorkspaceID string            `json:"workspaceId,omitempty"`
	Capability  string            `json:"capability"`
	Provider    string            `json:"provider,omitempty"`
	Target      string            `json:"target,omitempty"`
	Status      Status            `json:"status"`
	TaskID      string            `json:"taskId,omitempty"`
	Progress    int               `json:"progress,omitempty"`
	Detail      string            `json:"detail,omitempty"`
	Error       string            `json:"error,omitempty"`
	Result      any               `json:"result,omitempty"`
	Parameters  map[string]any    `json:"parameters,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	CreatedAt   time.Time         `json:"createdAt"`
	StartedAt   *time.Time        `json:"startedAt,omitempty"`
	CompletedAt *time.Time        `json:"completedAt,omitempty"`
}
