package resource

import (
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
	"time"
)

// ResourceSpec represents the desired state.
type ResourceSpec struct {
	common.Metadata
	EnvironmentID        common.EnvironmentID
	Type                 common.ResourceType
	Name                 valueobjects.Name
	TemplateID           common.TemplateID
	DesiredConfiguration string // Domain-agnostic payload
}

// RuntimeState represents the observed state.
type RuntimeState struct {
	ResourceID common.ResourceID
	Status     common.ResourceState
	Health     common.HealthState
	Metrics    map[string]interface{} // General purpose metrics container
	ObservedAt time.Time
}
