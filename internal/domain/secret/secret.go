package secret

import (
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
)

type Secret struct {
	common.Metadata
	EnvironmentID  common.EnvironmentID
	Name           valueobjects.Name
	Provider       string // e.g. "vault", "aws-secrets"
	Reference      valueobjects.Reference // External ID/path to the secret
	RotationPolicy string
	AccessPolicy   string
}
