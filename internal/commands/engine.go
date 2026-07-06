package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/Ajayvtl/devserver/internal/capabilities"
	"github.com/Ajayvtl/devserver/internal/events"
	rt "github.com/Ajayvtl/devserver/internal/runtime"
	"github.com/Ajayvtl/devserver/internal/state"
	"github.com/Ajayvtl/devserver/internal/tasks"
	"github.com/rs/zerolog"

	"github.com/Ajayvtl/devserver/internal/application/envcontext"
	"github.com/Ajayvtl/devserver/internal/domain/common")

// Engine is the command bus. It resolves capabilities and submits tasks.
type Engine struct {
	log      zerolog.Logger
	bus      events.Bus
	registry *capabilities.Registry
	resolver *capabilities.Resolver
	tasks    *tasks.Engine
	db       *state.StoreDB
	records  *store
	ctx      context.Context
	cancel   context.CancelFunc
	status   rt.Status
}

func NewEngine(log zerolog.Logger, bus events.Bus, registry *capabilities.Registry, taskEngine *tasks.Engine, db *state.StoreDB) *Engine {
	return &Engine{
		log:      log,
		bus:      bus,
		registry: registry,
		resolver: capabilities.NewResolver(registry),
		tasks:    taskEngine,
		db:       db,
		records:  newStore(),
		status:   rt.StatusStopped,
	}
}

func (e *Engine) Name() string { return "commands.Engine" }

func (e *Engine) Initialize(ctx context.Context) error {
	e.status = rt.StatusStarting
	return nil
}

func (e *Engine) Start(ctx context.Context) error {
	e.ctx, e.cancel = context.WithCancel(context.Background())
	e.status = rt.StatusRunning
	go e.runInternal(e.ctx)
	return nil
}

func (e *Engine) Stop(ctx context.Context) error {
	if e.cancel != nil {
		e.cancel()
	}
	e.status = rt.StatusStopped
	return nil
}

func (e *Engine) Status() rt.Status { return e.status }

func (e *Engine) Health() rt.Health { return rt.HealthHealthy }

func (e *Engine) Submit(ctx context.Context, cmd *Command) (*Record, error) {
	if cmd == nil {
		return nil, fmt.Errorf("command cannot be nil")
	}
	if cmd.ID == "" {
		cmd.ID = fmt.Sprintf("cmd-%d", time.Now().UnixNano())
	}
	if cmd.Name == "" {
		cmd.Name = cmd.Capability.String()
	}

	record := e.records.Create(cmd)
	e.bus.Publish(events.CommandSubmitted, events.CommandSubmittedEvent{
		CommandID:  record.ID,
		Name:       record.Name,
		Capability: record.Capability,
		Provider:   record.Provider,
	})

	if cmd.WorkspaceID != "" {
		switcher := envcontext.NewSwitcher(e.db)
		envID, _ := switcher.GetActiveEnvironment(ctx, common.WorkspaceID(cmd.WorkspaceID))
		ctx = switcher.WithActiveEnvironment(ctx, envID)
		if cmd.Metadata == nil {
			cmd.Metadata = make(map[string]string)
		}
		cmd.Metadata["env_id"] = string(envID)
	}

	switch cmd.Capability {
	case "auth.login":
		return e.runLogin(ctx, cmd, record)
	case "workspace.setup":
		return e.runSetup(ctx, cmd, record)
	case "project.save":
		return e.runProjectSave(ctx, cmd, record)
	case "workspace.index":
		return e.runWorkspaceIndex(ctx, cmd, record)
	}

	// Platform tasks bypass the module resolver
	if string(cmd.Capability) == "service.start" ||
		string(cmd.Capability) == "service.stop" ||
		string(cmd.Capability) == "service.restart" ||
		string(cmd.Capability) == "service.install" ||
		string(cmd.Capability) == "service.update" ||
		string(cmd.Capability) == "service.configure" ||
		string(cmd.Capability) == "service.remove" {
		return e.runServiceTask(ctx, cmd, record)
	}

	binding, err := e.resolver.Select(ctx, cmd.Capability, nil)
	if err != nil {
		record.Status = StatusFailed
		record.Error = err.Error()
		e.records.Update(record)
		e.bus.Publish(events.CommandFailed, events.CommandFailedEvent{
			CommandID:  record.ID,
			Name:       record.Name,
			Capability: record.Capability,
			Provider:   record.Provider,
			Error:      err.Error(),
		})
		return record, err
	}

	record.Provider = binding.Metadata.Provider
	e.records.Update(record)

	task := &tasks.Task{
		ID:          record.ID,
		Name:        record.Name,
		WorkspaceID: cmd.WorkspaceID,
		Type:        string(cmd.Capability),
		Priority:    tasks.TaskPriorityNormal,
		MaxRetries:  cmd.Retries,
		Timeout:     cmd.Timeout,
		Metadata:    cmd.Metadata,
		Payload:     cmd.Parameters,
	}

	taskRecord, err := e.tasks.Submit(task)
	if err != nil {
		record.Status = StatusFailed
		record.Error = err.Error()
		e.records.Update(record)
		e.bus.Publish(events.CommandFailed, events.CommandFailedEvent{
			CommandID:  record.ID,
			TaskID:     record.ID,
			Name:       record.Name,
			Capability: record.Capability,
			Provider:   record.Provider,
			Error:      err.Error(),
		})
		return record, err
	}

	record.TaskID = taskRecord.ID
	e.records.Update(record)
	return record, nil
}

func (e *Engine) runLogin(ctx context.Context, cmd *Command, record *Record) (*Record, error) {
	if e.db == nil {
		return e.fail(record, "platform database unavailable")
	}
	email, _ := cmd.Parameters["email"].(string)
	password, _ := cmd.Parameters["password"].(string)
	remember, _ := cmd.Parameters["remember"].(bool)

	now := time.Now().UTC()
	record.Status = StatusRunning
	record.StartedAt = &now
	e.records.Update(record)
	e.bus.Publish(events.CommandStarted, events.CommandStartedEvent{
		CommandID:  record.ID,
		Name:       record.Name,
		Capability: record.Capability,
		Provider:   record.Provider,
	})

	res, err := e.db.Auth(ctx, state.LoginRequestData{Email: email, Password: password, Remember: remember})
	if err != nil {
		return e.fail(record, err.Error())
	}

	record.Status = StatusCompleted
	record.Result = res
	now = time.Now().UTC()
	record.CompletedAt = &now
	e.records.Update(record)
	e.bus.Publish(events.CommandCompleted, events.CommandCompletedEvent{
		CommandID:  record.ID,
		Name:       record.Name,
		Capability: record.Capability,
		Provider:   record.Provider,
	})
	return record, nil
}

func (e *Engine) runSetup(ctx context.Context, cmd *Command, record *Record) (*Record, error) {
	if e.db == nil {
		return e.fail(record, "platform database unavailable")
	}

	now := time.Now().UTC()
	record.Status = StatusRunning
	record.StartedAt = &now
	e.records.Update(record)
	e.bus.Publish(events.CommandStarted, events.CommandStartedEvent{
		CommandID:  record.ID,
		Name:       record.Name,
		Capability: record.Capability,
		Provider:   record.Provider,
	})

	task := &tasks.Task{
		ID:          record.ID,
		Name:        record.Name,
		WorkspaceID: cmd.WorkspaceID,
		Type:        "workspace.setup",
		Priority:    tasks.TaskPriorityNormal,
		Metadata:    cmd.Metadata,
		Payload:     cmd.Parameters,
	}

	taskRecord, err := e.tasks.Submit(task)
	if err != nil {
		return e.fail(record, err.Error())
	}
	record.TaskID = taskRecord.ID
	record.Result = map[string]any{"taskId": taskRecord.ID, "status": "started"}
	e.records.Update(record)
	return record, nil
}

func (e *Engine) runProjectSave(ctx context.Context, cmd *Command, record *Record) (*Record, error) {
	if e.db == nil {
		return e.fail(record, "platform database unavailable")
	}

	now := time.Now().UTC()
	record.Status = StatusRunning
	record.StartedAt = &now
	e.records.Update(record)
	e.bus.Publish(events.CommandStarted, events.CommandStartedEvent{
		CommandID:  record.ID,
		Name:       record.Name,
		Capability: record.Capability,
		Provider:   record.Provider,
	})

	form := state.ProjectFormDataStore{}
	if v, ok := cmd.Parameters["name"].(string); ok {
		form.Name = v
	}
	if v, ok := cmd.Parameters["slug"].(string); ok {
		form.Slug = v
	}
	if v, ok := cmd.Parameters["description"].(string); ok {
		form.Description = v
	}
	if v, ok := cmd.Parameters["repository"].(string); ok {
		form.Repository = v
	}
	if v, ok := cmd.Parameters["branch"].(string); ok {
		form.Branch = v
	}
	if v, ok := cmd.Parameters["environment"].(string); ok {
		form.Environment = v
	}
	if v, ok := cmd.Parameters["owner"].(string); ok {
		form.Owner = v
	}
	if v, ok := cmd.Parameters["deployTarget"].(string); ok {
		form.DeployTarget = v
	}
	if v, ok := cmd.Parameters["healthCheck"].(string); ok {
		form.HealthCheck = v
	}
	if v, ok := cmd.Parameters["domains"].(string); ok {
		form.Domains = v
	}
	if v, ok := cmd.Parameters["autoDeploy"].(bool); ok {
		form.AutoDeploy = v
	}

	saved, err := e.db.SaveProject(ctx, form)
	if err != nil {
		return e.fail(record, err.Error())
	}

	record.Result = saved
	record.Progress = 100
	record.Detail = "Project saved"
	record.Status = StatusCompleted
	now = time.Now().UTC()
	record.CompletedAt = &now
	e.records.Update(record)
	e.bus.Publish(events.CommandCompleted, events.CommandCompletedEvent{
		CommandID:  record.ID,
		Name:       record.Name,
		Capability: record.Capability,
		Provider:   record.Provider,
	})
	return record, nil
}

func (e *Engine) runWorkspaceIndex(ctx context.Context, cmd *Command, record *Record) (*Record, error) {
	task := &tasks.Task{
		ID:          record.ID,
		Name:        record.Name,
		WorkspaceID: cmd.WorkspaceID,
		Type:        "workspace.index",
		Priority:    tasks.TaskPriorityNormal,
		Metadata:    cmd.Metadata,
		Payload:     cmd.Parameters,
	}

	taskRecord, err := e.tasks.Submit(task)
	if err != nil {
		return e.fail(record, err.Error())
	}
	record.TaskID = taskRecord.ID
	record.Result = map[string]any{"taskId": taskRecord.ID, "status": "started"}
	e.records.Update(record)
	return record, nil
}

func (e *Engine) runServiceTask(ctx context.Context, cmd *Command, record *Record) (*Record, error) {
	payload := cmd.Parameters
	if payload == nil {
		payload = make(map[string]any)
	}
	payload["target"] = cmd.Target

	task := &tasks.Task{
		ID:          record.ID,
		Name:        record.Name,
		WorkspaceID: cmd.WorkspaceID,
		Type:        string(cmd.Capability),
		Priority:    tasks.TaskPriorityHigh,
		Metadata:    cmd.Metadata,
		Payload:     payload,
	}

	taskRecord, err := e.tasks.Submit(task)
	if err != nil {
		return e.fail(record, err.Error())
	}
	record.TaskID = taskRecord.ID
	record.Result = map[string]any{"taskId": taskRecord.ID, "status": "started"}
	e.records.Update(record)
	return record, nil
}

func (e *Engine) fail(record *Record, message string) (*Record, error) {
	now := time.Now().UTC()
	record.Status = StatusFailed
	record.Error = message
	record.CompletedAt = &now
	e.records.Update(record)
	e.bus.Publish(events.CommandFailed, events.CommandFailedEvent{
		CommandID:  record.ID,
		Name:       record.Name,
		Capability: record.Capability,
		Provider:   record.Provider,
		Error:      message,
	})
	return record, fmt.Errorf("%s", message)
}

func (e *Engine) List() []*Record {
	return e.records.List()
}

func (e *Engine) Get(id string) (*Record, error) {
	record := e.records.Get(id)
	if record == nil {
		return nil, fmt.Errorf("command %q not found", id)
	}
	return record, nil
}

// Run keeps the command store synchronized with task lifecycle events.
func (e *Engine) runInternal(ctx context.Context) {
	topics := []events.EventType{
		events.TaskStarted,
		events.TaskProgress,
		events.TaskCompleted,
		events.TaskFailed,
		events.TaskCancelled,
		events.TaskRolledBack,
	}

	for _, topic := range topics {
		sub := e.bus.Subscribe(topic)
		go func(topic events.EventType, sub events.Subscriber) {
			for {
				select {
				case <-ctx.Done():
					return
				case evt, ok := <-sub:
					if !ok {
						return
					}
					e.syncFromTaskEvent(topic, evt.Payload)
				}
			}
		}(topic, sub)
	}
}

func (e *Engine) syncFromTaskEvent(topic events.EventType, payload any) {
	switch topic {
	case events.TaskStarted:
		if evt, ok := payload.(events.TaskStartedEvent); ok {
			if record := e.records.Get(evt.TaskID); record != nil {
				now := time.Now().UTC()
				record.Status = StatusRunning
				record.StartedAt = &now
				e.records.Update(record)
			}
		}
	case events.TaskProgress:
		if evt, ok := payload.(events.TaskProgressEvent); ok {
			if record := e.records.Get(evt.TaskID); record != nil {
				record.Progress = evt.Progress
				record.Detail = evt.Detail
				e.records.Update(record)
			}
		}
	case events.TaskCompleted:
		if evt, ok := payload.(events.TaskCompletedEvent); ok {
			if record := e.records.Get(evt.TaskID); record != nil {
				now := time.Now().UTC()
				record.Status = StatusCompleted
				record.Progress = 100
				record.CompletedAt = &now
				e.records.Update(record)
				e.bus.Publish(events.CommandCompleted, events.CommandCompletedEvent{
					CommandID:  record.ID,
					TaskID:     evt.TaskID,
					Name:       record.Name,
					Capability: record.Capability,
					Provider:   record.Provider,
				})
			}
		}
	case events.TaskFailed:
		if evt, ok := payload.(events.TaskFailedEvent); ok {
			if record := e.records.Get(evt.TaskID); record != nil {
				now := time.Now().UTC()
				record.Status = StatusFailed
				record.Error = evt.Error
				record.CompletedAt = &now
				e.records.Update(record)
				e.bus.Publish(events.CommandFailed, events.CommandFailedEvent{
					CommandID:  record.ID,
					TaskID:     evt.TaskID,
					Name:       record.Name,
					Capability: record.Capability,
					Provider:   record.Provider,
					Error:      evt.Error,
				})
			}
		}
	case events.TaskCancelled:
		if evt, ok := payload.(events.TaskCancelledEvent); ok {
			if record := e.records.Get(evt.TaskID); record != nil {
				now := time.Now().UTC()
				record.Status = StatusCancelled
				record.CompletedAt = &now
				e.records.Update(record)
				e.bus.Publish(events.CommandCancelled, events.CommandCancelledEvent{
					CommandID:  record.ID,
					TaskID:     evt.TaskID,
					Name:       record.Name,
					Capability: record.Capability,
					Provider:   record.Provider,
				})
			}
		}
	case events.TaskRolledBack:
		if evt, ok := payload.(events.TaskRolledBackEvent); ok {
			if record := e.records.Get(evt.TaskID); record != nil {
				now := time.Now().UTC()
				record.Status = StatusRolledBack
				record.CompletedAt = &now
				e.records.Update(record)
				e.bus.Publish(events.CommandRolledBack, events.CommandRolledBackEvent{
					CommandID:  record.ID,
					TaskID:     evt.TaskID,
					Name:       record.Name,
					Capability: record.Capability,
					Provider:   record.Provider,
				})
			}
		}
	}
}
