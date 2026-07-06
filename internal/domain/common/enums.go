package common

// EnvironmentType defines the purpose of the environment.
type EnvironmentType string

const (
	EnvironmentTypeLocal       EnvironmentType = "local"
	EnvironmentTypeDevelopment EnvironmentType = "development"
	EnvironmentTypeQA          EnvironmentType = "qa"
	EnvironmentTypeStaging     EnvironmentType = "staging"
	EnvironmentTypeProduction  EnvironmentType = "production"
)

// ExecutorType defines the mechanism for executing actions.
type ExecutorType string

const (
	ExecutorTypeLocal      ExecutorType = "local"
	ExecutorTypeSSH        ExecutorType = "ssh"
	ExecutorTypeDocker     ExecutorType = "docker"
	ExecutorTypeWSL        ExecutorType = "wsl"
	ExecutorTypeKubernetes ExecutorType = "kubernetes"
)

// ResourceType defines the classification of a resource.
type ResourceType string

const (
	ResourceTypeService  ResourceType = "service"
	ResourceTypeDatabase ResourceType = "database"
	ResourceTypeCache    ResourceType = "cache"
	ResourceTypeIngress  ResourceType = "ingress"
)

// ActionType defines common formalized intents.
type ActionType string

const (
	ActionTypeStart   ActionType = "start"
	ActionTypeStop    ActionType = "stop"
	ActionTypeRestart ActionType = "restart"
	ActionTypeInstall ActionType = "install"
	ActionTypeDetect  ActionType = "detect"
	ActionTypeBackup  ActionType = "backup"
	ActionTypeRestore ActionType = "restore"
)

// WorkflowState tracks the execution DAG status.
type WorkflowState string

const (
	WorkflowStatePending  WorkflowState = "pending"
	WorkflowStateRunning  WorkflowState = "running"
	WorkflowStateSuccess  WorkflowState = "success"
	WorkflowStateFailed   WorkflowState = "failed"
	WorkflowStateRollback WorkflowState = "rollback"
)

// ResourceState tracks the observed condition.
type ResourceState string

const (
	ResourceStateUnknown      ResourceState = "unknown"
	ResourceStateNotInstalled ResourceState = "not_installed"
	ResourceStateStopped      ResourceState = "stopped"
	ResourceStateStarting     ResourceState = "starting"
	ResourceStateRunning      ResourceState = "running"
	ResourceStateFailed       ResourceState = "failed"
	ResourceStateDegraded     ResourceState = "degraded"
)

// HealthState maps to resource health checks.
type HealthState string

const (
	HealthStateHealthy   HealthState = "healthy"
	HealthStateUnhealthy HealthState = "unhealthy"
	HealthStateCritical  HealthState = "critical"
	HealthStateUnknown   HealthState = "unknown"
)
