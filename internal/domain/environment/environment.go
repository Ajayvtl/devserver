package environment

import (
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
)

type Environment struct {
	common.Metadata
	WorkspaceID common.WorkspaceID
	Name        valueobjects.Name
	Type        common.EnvironmentType

	// Ownership refs
	Connections []common.ConnectionID
	Variables   []common.VariableID
	Secrets     []common.SecretID
	Policies    []common.PolicyID
	Resources   []common.ResourceID
	Workflows   []common.WorkflowID
	Templates   []common.TemplateID
}
