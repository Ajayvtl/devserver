// Package repository defines the canonical interfaces for data persistence.
// Note: Repositories NEVER publish events. Event publishing is strictly the responsibility
// of Application/Workflow services.
package repository

import (
	"context"

	"github.com/Ajayvtl/devserver/internal/domain/action"
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/domain/connection"
	"github.com/Ajayvtl/devserver/internal/domain/dependency"
	"github.com/Ajayvtl/devserver/internal/domain/environment"
	"github.com/Ajayvtl/devserver/internal/domain/event"
	"github.com/Ajayvtl/devserver/internal/domain/policy"
	"github.com/Ajayvtl/devserver/internal/domain/project"
	"github.com/Ajayvtl/devserver/internal/domain/resource"
	"github.com/Ajayvtl/devserver/internal/domain/secret"
	"github.com/Ajayvtl/devserver/internal/domain/template"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
	"github.com/Ajayvtl/devserver/internal/domain/variable"
	"github.com/Ajayvtl/devserver/internal/domain/workflow"
	"github.com/Ajayvtl/devserver/internal/domain/workspace"
)

// UnitOfWork abstracts database transactions for multi-entity operations.
type UnitOfWork interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type ProjectRepository interface {
	GetByID(ctx context.Context, id common.ProjectID) (*project.Project, error)
	List(ctx context.Context, opts common.ListOptions) ([]*project.Project, error)
	Create(ctx context.Context, entity *project.Project) error
	Update(ctx context.Context, entity *project.Project) error
	Delete(ctx context.Context, id common.ProjectID) error
	Exists(ctx context.Context, id common.ProjectID) (bool, error)
}

type WorkspaceRepository interface {
	GetByID(ctx context.Context, id common.WorkspaceID) (*workspace.Workspace, error)
	List(ctx context.Context, projectID common.ProjectID, opts common.ListOptions) ([]*workspace.Workspace, error)
	Create(ctx context.Context, entity *workspace.Workspace) error
	Update(ctx context.Context, entity *workspace.Workspace) error
	Delete(ctx context.Context, id common.WorkspaceID) error
	Exists(ctx context.Context, id common.WorkspaceID) (bool, error)
}

type EnvironmentRepository interface {
	GetByID(ctx context.Context, id common.EnvironmentID) (*environment.Environment, error)
	List(ctx context.Context, workspaceID common.WorkspaceID, opts common.ListOptions) ([]*environment.Environment, error)
	Create(ctx context.Context, entity *environment.Environment) error
	Update(ctx context.Context, entity *environment.Environment) error
	Delete(ctx context.Context, id common.EnvironmentID) error
	Exists(ctx context.Context, id common.EnvironmentID) (bool, error)
}

type TemplateRepository interface {
	GetByID(ctx context.Context, id common.TemplateID) (*template.Template, error)
	List(ctx context.Context, opts common.ListOptions) ([]*template.Template, error)
	Create(ctx context.Context, entity *template.Template) error
	Update(ctx context.Context, entity *template.Template) error
	Delete(ctx context.Context, id common.TemplateID) error
	Exists(ctx context.Context, id common.TemplateID) (bool, error)
}

type ResourceRepository interface {
	GetByID(ctx context.Context, id common.ResourceID) (*resource.ResourceSpec, error)
	List(ctx context.Context, envID common.EnvironmentID, opts common.ListOptions) ([]*resource.ResourceSpec, error)
	Create(ctx context.Context, entity *resource.ResourceSpec) error
	Update(ctx context.Context, entity *resource.ResourceSpec) error
	Delete(ctx context.Context, id common.ResourceID) error
	Exists(ctx context.Context, id common.ResourceID) (bool, error)
}

type RuntimeStateRepository interface {
	GetByID(ctx context.Context, id common.ResourceID) (*resource.RuntimeState, error)
	// List is intentionally omitted; RuntimeState is retrieved per-resource or aggregated via runtime metrics queries.
	// Save encompasses both Create and Update due to its ephemeral, rapid-update nature.
	Save(ctx context.Context, state *resource.RuntimeState) error
	Delete(ctx context.Context, id common.ResourceID) error
	Exists(ctx context.Context, id common.ResourceID) (bool, error)
}

type WorkflowRepository interface {
	GetByID(ctx context.Context, id common.WorkflowID) (*workflow.Workflow, error)
	List(ctx context.Context, envID common.EnvironmentID, opts common.ListOptions) ([]*workflow.Workflow, error)
	Create(ctx context.Context, entity *workflow.Workflow) error
	Update(ctx context.Context, entity *workflow.Workflow) error
	Delete(ctx context.Context, id common.WorkflowID) error
	Exists(ctx context.Context, id common.WorkflowID) (bool, error)
}

type ActionRepository interface {
	GetByID(ctx context.Context, id common.ActionID) (*action.Action, error)
	List(ctx context.Context, opts common.ListOptions) ([]*action.Action, error)
	Create(ctx context.Context, entity *action.Action) error
	Update(ctx context.Context, entity *action.Action) error
	Delete(ctx context.Context, id common.ActionID) error
	Exists(ctx context.Context, id common.ActionID) (bool, error)
}

type EventRepository interface {
	// Events are append-only; Update and Delete are intentionally omitted to maintain audit integrity.
	GetByID(ctx context.Context, id common.EventID) (*event.Event, error)
	List(ctx context.Context, opts common.ListOptions) ([]*event.Event, error)
	Create(ctx context.Context, entity *event.Event) error
	Exists(ctx context.Context, id common.EventID) (bool, error)
}

type VariableRepository interface {
	GetByID(ctx context.Context, id common.VariableID) (*variable.Variable, error)
	List(ctx context.Context, envID common.EnvironmentID, opts common.ListOptions) ([]*variable.Variable, error)
	Create(ctx context.Context, entity *variable.Variable) error
	Update(ctx context.Context, entity *variable.Variable) error
	Delete(ctx context.Context, id common.VariableID) error
	Exists(ctx context.Context, id common.VariableID) (bool, error)
}

type SecretRepository interface {
	GetByID(ctx context.Context, id common.SecretID) (*secret.Secret, error)
	List(ctx context.Context, envID common.EnvironmentID, opts common.ListOptions) ([]*secret.Secret, error)
	Create(ctx context.Context, entity *secret.Secret) error
	Update(ctx context.Context, entity *secret.Secret) error
	Delete(ctx context.Context, id common.SecretID) error
	Exists(ctx context.Context, id common.SecretID) (bool, error)
}

type ConnectionRepository interface {
	GetByID(ctx context.Context, id common.ConnectionID) (*connection.Connection, error)
	List(ctx context.Context, envID common.EnvironmentID, opts common.ListOptions) ([]*connection.Connection, error)
	Create(ctx context.Context, entity *connection.Connection) error
	Update(ctx context.Context, entity *connection.Connection) error
	Delete(ctx context.Context, id common.ConnectionID) error
	Exists(ctx context.Context, id common.ConnectionID) (bool, error)
}

type PolicyRepository interface {
	GetByID(ctx context.Context, id common.PolicyID) (*policy.Policy, error)
	List(ctx context.Context, envID common.EnvironmentID, opts common.ListOptions) ([]*policy.Policy, error)
	Create(ctx context.Context, entity *policy.Policy) error
	Update(ctx context.Context, entity *policy.Policy) error
	Delete(ctx context.Context, id common.PolicyID) error
	Exists(ctx context.Context, id common.PolicyID) (bool, error)
}

type DependencyRepository interface {
	ListBySource(ctx context.Context, source valueobjects.Reference, opts common.ListOptions) ([]*dependency.Dependency, error)
	ListByTarget(ctx context.Context, target valueobjects.Reference, opts common.ListOptions) ([]*dependency.Dependency, error)
	Create(ctx context.Context, entity *dependency.Dependency) error
	// Dependencies are relational mapping tables; they lack a single ID and updates are usually Delete + Create.
	Delete(ctx context.Context, source, target valueobjects.Reference) error
}
