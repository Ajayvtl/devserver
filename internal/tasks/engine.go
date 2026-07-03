package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/Ajayvtl/devserver/internal/events"
	"github.com/rs/zerolog"
)

// Engine is the central task execution coordinator.
// It owns the queue, worker pool, and task store.
type Engine struct {
	log     zerolog.Logger
	bus     events.Bus
	queue   *queue
	store   *store
	workers []*worker
	exec    Executor
	size    int
}

// EngineConfig controls engine behavior.
type EngineConfig struct {
	Workers  int // Number of concurrent workers. Default: 4.
	QueueMax int // Maximum queue depth. 0 = unbounded.
}

// NewEngine creates a task engine with the given configuration.
func NewEngine(log zerolog.Logger, bus events.Bus, exec Executor, cfg EngineConfig) *Engine {
	if cfg.Workers <= 0 {
		cfg.Workers = 4
	}

	q := newQueue(cfg.QueueMax)
	s := newStore()

	if exec == nil {
		exec = NewDirectExecutor()
	}

	workers := make([]*worker, cfg.Workers)
	for i := 0; i < cfg.Workers; i++ {
		workers[i] = newWorker(i, log, bus, exec, q, s)
	}

	return &Engine{
		log:     log,
		bus:     bus,
		queue:   q,
		store:   s,
		workers: workers,
		exec:    exec,
		size:    cfg.Workers,
	}
}

// Run starts all workers. Blocks until ctx is cancelled.
func (e *Engine) Run(ctx context.Context) {
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

// Submit enqueues a task definition for execution.
// Returns the task record, or an error if the queue is full.
func (e *Engine) Submit(def *Definition) (*Record, error) {
	if def.ID == "" {
		def.ID = fmt.Sprintf("task-%d", time.Now().UnixNano())
	}
	if def.Retry.MaxRetries == 0 && def.Retries > 0 {
		def.Retry.MaxRetries = def.Retries
	}
	if def.Retry.Backoff <= 0 {
		def.Retry.Backoff = defaultRetryBackoff
	}

	record := e.store.Create(def)

	if !e.queue.Enqueue(def) {
		record.Status = StatusFailed
		record.Error = "queue is full"
		e.store.Update(record)
		e.bus.Publish(events.TaskFailed, events.TaskFailedEvent{
			TaskID: def.ID,
			Name:   def.Name,
			Error:  "queue is full",
		})
		return record, fmt.Errorf("task queue is full")
	}

	publishTaskQueued(e.bus, record)
	e.log.Debug().Str("task", def.ID).Str("name", def.Name).Msg("task submitted")
	return record, nil
}

// Cancel stops a running or queued task.
func (e *Engine) Cancel(id string) error {
	// Try removing from queue first.
	if e.queue.Remove(id) {
		record := e.store.Get(id)
		if record != nil {
			record.markCancelled("Cancelled while queued")
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
				record.markCancelled("Cancelled while running")
				e.store.Update(record)
				publishTaskCancelled(e.bus, record)
			}
			return nil
		}
	}

	return fmt.Errorf("task %q not found in queue or running", id)
}

// Get returns the current state of a task.
func (e *Engine) Get(id string) (*Record, error) {
	r := e.store.Get(id)
	if r == nil {
		return nil, fmt.Errorf("task %q not found", id)
	}
	return r, nil
}

// List returns all task records.
func (e *Engine) List() []*Record {
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
