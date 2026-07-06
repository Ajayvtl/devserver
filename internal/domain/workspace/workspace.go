package workspace

import (
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
)

type Workspace struct {
	common.Metadata
	ProjectID common.ProjectID
	Path      valueobjects.Path
}
