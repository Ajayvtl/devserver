package application

import (
	"context"

	"github.com/Ajayvtl/devserver/internal/domain/common"
)

// EnvironmentService manages the high-level lifecycle of environments.
type EnvironmentService interface {
	Clone(ctx context.Context, sourceEnvID, targetEnvID common.EnvironmentID) error
	Archive(ctx context.Context, envID common.EnvironmentID) error
}

// ResourceService manages the high-level lifecycle of resources.
type ResourceService interface {
	Reconcile(ctx context.Context, resID common.ResourceID) error
}
