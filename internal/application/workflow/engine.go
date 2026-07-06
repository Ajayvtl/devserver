package workflow

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

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
	StatusPending    = "pending"
	StatusRunning    = "running"
	StatusPaused     = "paused"
	StatusCancelled  = "cancelled"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
	StatusRolledBack = "rolled_back"
)

type NodeState struct {
	Node    *domainWorkflow.WorkflowNode
	Status  string
	Error   error
	Retries int
}

type WorkflowState struct {
	Workflow *domainWorkflow.Workflow
	Status   string
	Nodes    map[string]*NodeState
	Cancel   context.CancelFunc
	ctx      context.Context
	pauseCh  chan struct{}
}

type EngineImpl struct {
	dispatcher Dispatcher
	registry   contracts.ExecutorRegistry
	bus        events.Bus

	mu       sync.RWMutex
	states   map[common.WorkflowID]*WorkflowState
	nodeToWf map[string]common.WorkflowID
}

func NewEngine(dispatcher Dispatcher, registry contracts.ExecutorRegistry, bus events.Bus) Engine {
	return &EngineImpl{
		dispatcher: dispatcher,
		registry:   registry,
		bus:        bus,
		states:     make(map[common.WorkflowID]*WorkflowState),
		nodeToWf:   make(map[string]common.WorkflowID),
	}
}

func (e *EngineImpl) CreateExecutionPlan(ctx context.Context, wf *domainWorkflow.Workflow) ([]*domainWorkflow.WorkflowNode, error) {
	// Simple topological sort to create plan
	inDegree := make(map[string]int)
	graph := make(map[string][]string)

	e.mu.Lock()
	state, exists := e.states[common.WorkflowID(wf.ID)]
	if !exists {
		state = &WorkflowState{
			Workflow: wf,
			Status:   StatusPending,
			Nodes:    make(map[string]*NodeState),
			pauseCh:  make(chan struct{}),
		}
		e.states[common.WorkflowID(wf.ID)] = state
	}
	for _, n := range wf.Nodes {
		inDegree[n.ID] = 0
		graph[n.ID] = []string{}
		e.nodeToWf[n.ID] = common.WorkflowID(wf.ID)
	}
	e.mu.Unlock()

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

func (e *EngineImpl) waitIfPaused(ctx context.Context, nodeID string) error {
	e.mu.RLock()
	wfID, ok := e.nodeToWf[nodeID]
	if !ok {
		e.mu.RUnlock()
		return nil
	}
	state, ok := e.states[wfID]
	e.mu.RUnlock()

	if !ok {
		return nil
	}

	for {
		e.mu.RLock()
		status := state.Status
		pauseCh := state.pauseCh
		e.mu.RUnlock()

		if status != StatusPaused {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-pauseCh:
			// Woken up, loop and check status again
		}
	}
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

	backoff := time.Second

	for i := 0; i <= retries; i++ {
		// Check cancellation before execution
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if waitErr := e.waitIfPaused(ctx, node.ID); waitErr != nil {
			return waitErr
		}

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

		if i < retries {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
				backoff *= 2
			}
		}
	}

	if err != nil && node.Rollback != "" {
		rollbackNode := &domainWorkflow.WorkflowNode{
			ID:       node.ID + "-rollback",
			ActionID: node.Rollback,
		}

		rollbackErr := e.dispatcher.Dispatch(ctx, rollbackNode)
		if rollbackErr != nil {
			err = fmt.Errorf("node failed: %v, rollback also failed: %v", err, rollbackErr)
		}

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
	close(state.pauseCh)
	state.pauseCh = make(chan struct{})
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
