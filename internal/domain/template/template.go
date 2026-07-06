package template

import (
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
)

type Template struct {
	common.Metadata
	Name    valueobjects.Name
	Version valueobjects.Version
	Schema  string // Domain-agnostic definition payload
}
