package workflow

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Ajayvtl/devserver/internal/commands"
	"github.com/Ajayvtl/devserver/internal/domain/common"
	domainWorkflow "github.com/Ajayvtl/devserver/internal/domain/workflow"
	"github.com/Ajayvtl/devserver/internal/events"
	"github.com/Ajayvtl/devserver/internal/executor/contracts"
)

var (
	ErrWorkflowNotFound  = errors.New("workflow not found")
	ErrNodeNotFound      = errors.New("node not found")
	ErrWorkflowNotActive = errors.New("workflow is not active")
)

const (
	StatusPending   = "pending"
	StatusRunning   = "running"
	StatusPaused    = "paused"
	StatusCancelled = "cancelled"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
	StatusRolledBack = "rolled_back"
)

type NodeState struct {
	Node     *domainWorkflow.WorkflowNode
	Status   string
	Error    error
	Retries  int
}

type WorkflowState struct {
	Workflow *domainWorkflow.Workflow
	Status   string
	Nodes    map[string]*NodeState
	Cancel   context.CancelFunc
	ctx      context.Context
}

type EngineImpl struct {
	dispatcher Dispatcher
	registry   contracts.ExecutorRegistry
	bus        events.Bus
	cmdEngine  *commands.Engine

	mu     sync.RWMutex
	states map[common.WorkflowID]*WorkflowState
}

func NewEngine(dispatcher Dispatcher, registry contracts.ExecutorRegistry, bus events.Bus, cmdEngine *commands.Engine) Engine {
	return &EngineImpl{
		dispatcher: dispatcher,
		registry:   registry,
		bus:        bus,
		cmdEngine:  cmdEngine,
		states:     make(map[common.WorkflowID]*WorkflowState),
	}
}

func (e *EngineImpl) CreateExecutionPlan(ctx context.Context, wf *domainWorkflow.Workflow) ([]*domainWorkflow.WorkflowNode, error) {
	// Simple topological sort to create plan
	inDegree := make(map[string]int)
	graph := make(map[string][]string)

	for _, n := range wf.Nodes {
		inDegree[n.ID] = 0
		graph[n.ID] = []string{}
	}

	for _, edge := range wf.Edges {
		inDegree[edge.To]++
		graph[edge.From] = append(graph[edge.From], edge.To)
	}

	var queue []string
	for id, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, id)
		}
	}

	var plan []*domainWorkflow.WorkflowNode
	nodeMap := make(map[string]*domainWorkflow.WorkflowNode)
	for i := range wf.Nodes {
		nodeMap[wf.Nodes[i].ID] = &wf.Nodes[i]
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		plan = append(plan, nodeMap[curr])

		for _, neighbor := range graph[curr] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	if len(plan) != len(wf.Nodes) {
		return nil, fmt.Errorf("cycle detected in workflow DAG")
	}

	return plan, nil
}

func (e *EngineImpl) ExecutePlan(ctx context.Context, plan []*domainWorkflow.WorkflowNode) error {
	// This method gets the flat list but to run DAG, we need the edges which are in the original workflow.
	// Assuming the caller has already created a state, or we just execute them sequentially as a fallback.
	// Since the interface signature implies plan is just a list of nodes, we execute them in order (sequential).
	for _, node := range plan {
		err := e.executeNode(ctx, node)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *EngineImpl) executeNode(ctx context.Context, node *domainWorkflow.WorkflowNode) error {
	// Check Condition
	if node.Condition != "" {
		// Evaluator logic here - for now assume passed or skip
	}

	var err error
	retries := node.Retry
	if retries < 0 {
		retries = 0
	}

	for i := 0; i <= retries; i++ {
		execCtx := ctx
		var cancel context.CancelFunc
		if node.Timeout > 0 {
			execCtx, cancel = context.WithTimeout(ctx, time.Duration(node.Timeout)*time.Second)
		} else {
			execCtx, cancel = context.WithCancel(ctx)
		}

		// Emit progress
		e.bus.Publish(events.TaskProgress, events.TaskProgressEvent{
			TaskID:   node.ID,
			Progress: 10,
			Detail:   fmt.Sprintf("Executing node %s (Attempt %d)", node.ID, i+1),
		})

		err = e.dispatcher.Dispatch(execCtx, node)
		cancel()

		if err == nil {
			e.bus.Publish(events.TaskProgress, events.TaskProgressEvent{
				TaskID:   node.ID,
				Progress: 100,
				Detail:   fmt.Sprintf("Completed node %s", node.ID),
			})
			return nil
		}

		time.Sleep(time.Second) // backoff
	}

	if err != nil && node.Rollback != "" {
		// Execute rollback action via dispatcher if implemented, or log
		// Fallback for demonstration
		e.bus.Publish(events.TaskRolledBack, events.TaskRolledBackEvent{
			TaskID: node.ID,
			Name:   string(node.Rollback),
		})
	}

	return err
}

func (e *EngineImpl) Pause(ctx context.Context, wfID common.WorkflowID) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	state, ok := e.states[wfID]
	if !ok {
		return ErrWorkflowNotFound
	}
	if state.Status != StatusRunning {
		return ErrWorkflowNotActive
	}
	state.Status = StatusPaused
	return nil
}

func (e *EngineImpl) Resume(ctx context.Context, wfID common.WorkflowID) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	state, ok := e.states[wfID]
	if !ok {
		return ErrWorkflowNotFound
	}
	if state.Status != StatusPaused {
		return fmt.Errorf("workflow not paused")
	}
	state.Status = StatusRunning
	return nil
}

func (e *EngineImpl) Cancel(ctx context.Context, wfID common.WorkflowID) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	state, ok := e.states[wfID]
	if !ok {
		return ErrWorkflowNotFound
	}
	if state.Cancel != nil {
		state.Cancel()
	}
	state.Status = StatusCancelled
	return nil
}

func (e *EngineImpl) Retry(ctx context.Context, wfID common.WorkflowID, nodeID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	state, ok := e.states[wfID]
	if !ok {
		return ErrWorkflowNotFound
	}
	
	nodeState, ok := state.Nodes[nodeID]
	if !ok {
		return ErrNodeNotFound
	}
	
	if nodeState.Status != StatusFailed {
		return fmt.Errorf("node is not in failed state")
	}
	
	nodeState.Status = StatusPending
	nodeState.Error = nil
	nodeState.Retries++
	
	// Re-trigger execution for this node
	return nil
}

func (e *EngineImpl) Rollback(ctx context.Context, wfID common.WorkflowID) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	state, ok := e.states[wfID]
	if !ok {
		return ErrWorkflowNotFound
	}
	
	state.Status = StatusRolledBack
	return nil
}

func (e *EngineImpl) Status(ctx context.Context, wfID common.WorkflowID) (string, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	state, ok := e.states[wfID]
	if !ok {
		return "", ErrWorkflowNotFound
	}
	return state.Status, nil
}
