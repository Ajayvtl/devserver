package dependency

import (
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
)

type DependencyType string

const (
	DependencyTypeSoft    DependencyType = "soft"
	DependencyTypeHard    DependencyType = "hard"
	DependencyTypeRuntime DependencyType = "runtime"
	DependencyTypeBuild   DependencyType = "build"
)

type Dependency struct {
	common.Metadata
	From      valueobjects.Reference
	To        valueobjects.Reference
	Type      DependencyType
	Condition string
}
