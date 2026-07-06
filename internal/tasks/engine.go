package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/Ajayvtl/devserver/internal/events"
	rt "github.com/Ajayvtl/devserver/internal/runtime"
	"github.com/rs/zerolog"
)

// Engine is the central task execution coordinator.
// It owns the queue, worker pool, task store, registry, and runtime.
type Engine struct {
	log      zerolog.Logger
	bus      events.Bus
	queue    *queue
	store    *store
	workers  []*worker
	registry *Registry
	runtime  *Runtime
	size     int
	ctx      context.Context
	cancel   context.CancelFunc
	status   rt.Status
}

// EngineConfig controls engine behavior.
type EngineConfig struct {
	Workers  int // Number of concurrent workers. Default: 4.
	QueueMax int // Maximum queue depth. 0 = unbounded.
}

// NewEngine creates a task engine with the given configuration.
func NewEngine(log zerolog.Logger, bus events.Bus, registry *Registry, runtime *Runtime, cfg EngineConfig) *Engine {
	if cfg.Workers <= 0 {
		cfg.Workers = 4
	}

	q := newQueue(cfg.QueueMax)
	s := newStore()

	if registry == nil {
		registry = NewRegistry()
	}

	workers := make([]*worker, cfg.Workers)
	for i := 0; i < cfg.Workers; i++ {
		workers[i] = newWorker(i, log, bus, registry, runtime, q, s)
	}

	return &Engine{
		log:      log,
		bus:      bus,
		queue:    q,
		store:    s,
		workers:  workers,
		registry: registry,
		runtime:  runtime,
		size:     cfg.Workers,
		status:   rt.StatusStopped,
	}
}

func (e *Engine) Name() string { return "tasks.Engine" }

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

// runInternal starts all workers. Blocks until ctx is cancelled.
func (e *Engine) runInternal(ctx context.Context) {
	e.log.Info().Int("workers", e.size).Msg("task engine started")

	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		e.queue.Close()
		close(done)
	}()

	for _, w := range e.workers {
		go w.run(ctx)
	}

	<-done
	e.log.Info().Msg("task engine stopped")
}

// Submit enqueues a task for execution.
// Returns the task, or an error if the queue is full.
func (e *Engine) Submit(task *Task) (*Task, error) {
	if task.ID == "" {
		task.ID = fmt.Sprintf("task-%d", time.Now().UnixNano())
	}
	if task.Status == "" {
		task.Status = TaskQueued
	}
	if task.StartedAt.IsZero() {
		task.StartedAt = time.Now().UTC()
	}

	record := e.store.Create(task)

	if !e.queue.Enqueue(task) {
		record.Status = TaskFailed
		record.Error = "queue is full"
		e.store.Update(record)
		e.bus.Publish(events.TaskFailed, events.TaskFailedEvent{
			TaskID: task.ID,
			Name:   task.Name,
			Error:  "queue is full",
		})
		return record, fmt.Errorf("task queue is full")
	}

	publishTaskQueued(e.bus, record)
	e.log.Debug().Str("task", task.ID).Str("name", task.Name).Msg("task submitted")
	return record, nil
}

// Cancel stops a running or queued task.
func (e *Engine) Cancel(id string) error {
	// Try removing from queue first.
	if e.queue.Remove(id) {
		record := e.store.Get(id)
		if record != nil {
			record.Status = TaskCancelled
			e.store.Update(record)
			publishTaskCancelled(e.bus, record)
		}
		return nil
	}

	// Try cancelling a running task.
	for _, w := range e.workers {
		if w.cancelTask(id) {
			record := e.store.Get(id)
			if record != nil {
				record.Status = TaskCancelled
				e.store.Update(record)
				publishTaskCancelled(e.bus, record)
			}
			return nil
		}
	}

	return fmt.Errorf("task %q not found in queue or running", id)
}

// Get returns the current state of a task.
func (e *Engine) Get(id string) (*Task, error) {
	r := e.store.Get(id)
	if r == nil {
		return nil, fmt.Errorf("task %q not found", id)
	}
	return r, nil
}

// List returns all task records.
func (e *Engine) List() []*Task {
	return e.store.List()
}

// QueueLen returns the number of tasks waiting in the queue.
func (e *Engine) QueueLen() int {
	return e.queue.Len()
}

// Pause stops workers from draining the queue.
func (e *Engine) Pause() {
	e.queue.Pause()
}

// Resume resumes queue draining.
func (e *Engine) Resume() {
	e.queue.Resume()
}

// IsPaused reports whether the queue is paused.
func (e *Engine) IsPaused() bool {
	return e.queue.IsPaused()
}
