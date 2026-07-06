package event

import (
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
	"time"
)

type Event struct {
	common.Metadata
	Type             string
	Source           valueobjects.Reference
	Target           valueobjects.Reference
	PayloadReference valueobjects.Reference // Pointer to external payload storage, avoiding map[string]any
	OccurredAt       time.Time
}
