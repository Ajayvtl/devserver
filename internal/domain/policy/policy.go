package policy

import (
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
)

type Policy struct {
	common.Metadata
	EnvironmentID common.EnvironmentID
	Name          valueobjects.Name
	RuleSet       string // Domain-agnostic definition payload
}
