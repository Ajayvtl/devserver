package workflow

import (
	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
)

type TriggerType string

const (
	TriggerManual   TriggerType = "manual"
	TriggerSchedule TriggerType = "schedule"
	TriggerWebhook  TriggerType = "webhook"
	TriggerEvent    TriggerType = "event"
	TriggerAPI      TriggerType = "api"
)

type Workflow struct {
	common.Metadata
	EnvironmentID common.EnvironmentID
	Name          valueobjects.Name
	Trigger       TriggerType
	Nodes         []WorkflowNode
	Edges         []WorkflowEdge
}

type WorkflowNode struct {
	ID        string
	ActionID  common.ActionID
	Retry     int
	Timeout   int // Seconds
	Rollback  common.ActionID
	Condition string
}

type WorkflowEdge struct {
	From string
	To   string
}
