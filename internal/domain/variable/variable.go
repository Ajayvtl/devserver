package variable

import (
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
)

type Variable struct {
	common.Metadata
	EnvironmentID common.EnvironmentID
	Name          valueobjects.Name
	Value         string
}
