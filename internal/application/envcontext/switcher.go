package envcontext

import (
	"context"

	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/state"
)

type Switcher interface {
	// SetActiveEnvironment updates the active environment for a workspace in persistence.
	SetActiveEnvironment(ctx context.Context, workspaceID common.WorkspaceID, envID common.EnvironmentID) error

	// GetActiveEnvironment returns the active environment for a workspace from persistence.
	GetActiveEnvironment(ctx context.Context, workspaceID common.WorkspaceID) (common.EnvironmentID, error)

	// WithActiveEnvironment adds the given environment to the context for propagation.
	WithActiveEnvironment(ctx context.Context, envID common.EnvironmentID) context.Context

	// FromContext extracts the active environment from the context.
	FromContext(ctx context.Context) (common.EnvironmentID, bool)
}

type contextKey string

const activeEnvKey contextKey = "active_environment"

type storeSwitcher struct {
	db *state.StoreDB
}

func NewSwitcher(db *state.StoreDB) Switcher {
	return &storeSwitcher{db: db}
}

func (s *storeSwitcher) SetActiveEnvironment(ctx context.Context, workspaceID common.WorkspaceID, envID common.EnvironmentID) error {
	// In the current DB schema, the active environment is stored in the projects table.
	// Since we are not changing the DB schema, we update the existing project record.
	return s.db.UpdateProjectEnvironment(ctx, string(workspaceID), string(envID))
}

func (s *storeSwitcher) GetActiveEnvironment(ctx context.Context, workspaceID common.WorkspaceID) (common.EnvironmentID, error) {
	// Fallback to "development" if not found or DB unavailable
	if s.db == nil {
		return common.EnvironmentID("development"), nil
	}

	project, err := s.db.ProjectDetail(ctx, string(workspaceID))
	if err != nil {
		return common.EnvironmentID("development"), nil
	}

	if project.Project.Environment != "" {
		return common.EnvironmentID(project.Project.Environment), nil
	}

	return common.EnvironmentID("development"), nil
}

func (s *storeSwitcher) WithActiveEnvironment(ctx context.Context, envID common.EnvironmentID) context.Context {
	return context.WithValue(ctx, activeEnvKey, envID)
}

func (s *storeSwitcher) FromContext(ctx context.Context) (common.EnvironmentID, bool) {
	envID, ok := ctx.Value(activeEnvKey).(common.EnvironmentID)
	return envID, ok
}
