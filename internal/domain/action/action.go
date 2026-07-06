package action

import (
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
)

type Action struct {
	common.Metadata
	Type      common.ActionType
	Target    valueobjects.Reference
	Arguments []string
}
