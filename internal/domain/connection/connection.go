package connection

import (
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
)

type ConnectionType string

const (
	ConnectionTypeTransport   ConnectionType = "transport"
	ConnectionTypeIntegration ConnectionType = "integration"
)

type Connection struct {
	common.Metadata
	EnvironmentID common.EnvironmentID
	Name          valueobjects.Name
	Type          ConnectionType
	Protocol      string // e.g. SSH, HTTP, GitHub, AWS
	Endpoint      valueobjects.Endpoint
	SecretRef     common.SecretID
}
