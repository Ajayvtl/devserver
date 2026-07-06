package application

import (
	"context"

	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/domain/resource"
)

// Use cases define the business processes of the system.

type CreateEnvironmentUseCase interface {
	Execute(ctx context.Context, name string, envType common.EnvironmentType) (common.EnvironmentID, error)
}

type ProvisionResourceUseCase interface {
	Execute(ctx context.Context, envID common.EnvironmentID, spec *resource.ResourceSpec) (common.ResourceID, error)
}

type StartServiceUseCase interface {
	Execute(ctx context.Context, resID common.ResourceID) error
}
