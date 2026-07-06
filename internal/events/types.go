package events

type EventType string

const (
	WorkspaceChanged     EventType = "workspace:changed"
	WorkspaceIndexed     EventType = "workspace:indexed"
	WorkspaceActivated   EventType = "workspace:activated"
	WorkspaceDeactivated EventType = "workspace:deactivated"

	TaskStarted    EventType = "task:started"
	TaskProgress   EventType = "task:progress"
	TaskQueued     EventType = "task:queued"
	TaskCompleted  EventType = "task:completed"
	TaskFailed     EventType = "task:failed"
	TaskCancelled  EventType = "task:cancelled"
	TaskRolledBack EventType = "task:rolled_back"

	CommandSubmitted  EventType = "command:submitted"
	CommandStarted    EventType = "command:started"
	CommandCompleted  EventType = "command:completed"
	CommandFailed     EventType = "command:failed"
	CommandCancelled  EventType = "command:cancelled"
	CommandRolledBack EventType = "command:rolled_back"

	GitChanged  EventType = "git:changed"
	FileChanged EventType = "file:changed"

	DeploymentStarted  EventType = "deployment:started"
	DeploymentFinished EventType = "deployment:finished"

	PluginLoaded EventType = "plugin:loaded"
	PluginFailed EventType = "plugin:failed"

	ContextGenerated EventType = "context:generated"
	ContextUpdated   EventType = "context:updated"
)

// Typed event payloads. Every Publish call should use one of these
// instead of raw strings or untyped maps.

type WorkspaceChangedEvent struct {
	WorkspaceID string
}

type WorkspaceIndexedEvent struct {
	WorkspaceID  string
	ChangedFiles int
}

type WorkspaceActivatedEvent struct {
	WorkspaceID string
	Root        string
}

type WorkspaceDeactivatedEvent struct {
	WorkspaceID string
}

type TaskStartedEvent struct {
	TaskID string
	Name   string
}

type TaskQueuedEvent struct {
	TaskID      string
	Name        string
	WorkspaceID string
	Detail      string
}

type TaskProgressEvent struct {
	TaskID   string
	Progress int
	Detail   string
}

type TaskCompletedEvent struct {
	TaskID string
	Name   string
}

type TaskFailedEvent struct {
	TaskID string
	Name   string
	Error  string
}

type TaskCancelledEvent struct {
	TaskID string
	Name   string
}

type TaskRolledBackEvent struct {
	TaskID string
	Name   string
}

type CommandSubmittedEvent struct {
	CommandID  string
	Name       string
	Capability string
	Provider   string
}

type CommandStartedEvent struct {
	CommandID  string
	TaskID     string
	Name       string
	Capability string
	Provider   string
}

type CommandCompletedEvent struct {
	CommandID  string
	TaskID     string
	Name       string
	Capability string
	Provider   string
}

type CommandFailedEvent struct {
	CommandID  string
	TaskID     string
	Name       string
	Capability string
	Provider   string
	Error      string
}

type CommandCancelledEvent struct {
	CommandID  string
	TaskID     string
	Name       string
	Capability string
	Provider   string
}

type CommandRolledBackEvent struct {
	CommandID  string
	TaskID     string
	Name       string
	Capability string
	Provider   string
}

type FileChangedEvent struct {
	WorkspaceID string
	Path        string
}

type GitChangedEvent struct {
	WorkspaceID string
	Branch      string
}

type ContextGeneratedEvent struct {
	WorkspaceID string
}
