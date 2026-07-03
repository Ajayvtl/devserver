package tasks

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Ajayvtl/devserver/internal/events"
	"github.com/rs/zerolog"
)

// worker pulls tasks from the queue and executes them.
type worker struct {
	id     int
	log    zerolog.Logger
	bus    events.Bus
	exec   Executor
	queue  *queue
	store  *store
	cancel map[string]context.CancelFunc
	mu     sync.Mutex
}

func newWorker(id int, log zerolog.Logger, bus events.Bus, exec Executor, q *queue, s *store) *worker {
	return &worker{
		id:     id,
		log:    log.With().Int("worker", id).Logger(),
		bus:    bus,
		exec:   exec,
		queue:  q,
		store:  s,
		cancel: make(map[string]context.CancelFunc),
	}
}

// run processes tasks until ctx is cancelled.
func (w *worker) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-w.queue.Wait():
			if !ok {
				return
			}
			w.drain(ctx)
		case _, ok := <-w.queue.ResumeWait():
			if !ok {
				return
			}
			w.drain(ctx)
		}
	}
}

// drain processes all available tasks in the queue.
func (w *worker) drain(ctx context.Context) {
	for {
		if w.queue.IsPaused() {
			return
		}
		def := w.queue.Dequeue()
		if def == nil {
			return
		}
		w.execute(ctx, def)
	}
}

func (w *worker) execute(parentCtx context.Context, def *Definition) {
	record := w.store.Get(def.ID)
	if record == nil {
		return
	}

	// Set up cancellable context with optional timeout.
	taskCtx, cancel := context.WithCancel(parentCtx)
	if def.Timeout > 0 {
		taskCtx, cancel = context.WithTimeout(parentCtx, def.Timeout)
	}
	defer cancel()

	w.mu.Lock()
	w.cancel[def.ID] = cancel
	w.mu.Unlock()
	defer func() {
		w.mu.Lock()
		delete(w.cancel, def.ID)
		w.mu.Unlock()
	}()

	// Transition to running.
	now := time.Now().UTC()
	record.Status = StatusRunning
	record.StartedAt = &now
	record.Attempt++
	w.store.Update(record)
	w.bus.Publish(events.TaskStarted, events.TaskStartedEvent{
		TaskID: def.ID,
		Name:   def.Name,
	})
	w.log.Info().Str("task", def.ID).Str("name", def.Name).Msg("task started")

	// Build the task context for progress reporting.
	done := taskCtx.Done()
	tctx := &Context{
		record: record,
		progress: func(pct int, detail string) {
			record.Progress = pct
			record.Detail = detail
			w.store.Update(record)
			w.bus.Publish(events.TaskProgress, events.TaskProgressEvent{
				TaskID:   def.ID,
				Progress: pct,
				Detail:   detail,
			})
		},
		done: done,
	}

	// Execute the work.
	err := w.exec.Run(taskCtx, def, tctx)

	if err != nil {
		var execErr *ExecutionError
		if errors.As(err, &execErr) && execErr.RolledBack {
			w.bus.Publish(events.TaskFailed, events.TaskFailedEvent{
				TaskID: def.ID,
				Name:   def.Name,
				Error:  err.Error(),
			})
			record.markRolledBack("Rolled back after failure")
			w.store.Update(record)
			publishTaskRolledBack(w.bus, record)
			w.log.Warn().Err(err).Str("task", def.ID).Msg("task rolled back")
		} else if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			record.markCancelled("Cancelled")
			w.store.Update(record)
			publishTaskCancelled(w.bus, record)
			w.log.Warn().Err(err).Str("task", def.ID).Msg("task cancelled")
			return
		} else {
			w.log.Error().Err(err).Str("task", def.ID).Msg("task failed")
			record.markFailed(err.Error())
			w.bus.Publish(events.TaskFailed, events.TaskFailedEvent{
				TaskID: def.ID,
				Name:   def.Name,
				Error:  err.Error(),
			})
			w.store.Update(record)
		}

		// Retry with backoff if retries remain.
		if record.Attempt <= maxRetries(def) && maxRetries(def) > 0 {
			backoff := retryBackoff(def, record.Attempt)
			record.Status = StatusQueued
			record.Detail = "Retrying after " + backoff.String()
			record.CompletedAt = nil
			w.store.Update(record)
			w.log.Info().Str("task", def.ID).Dur("backoff", backoff).Msg("scheduling retry")

			go func() {
				select {
				case <-parentCtx.Done():
					return
				case <-time.After(backoff):
					w.queue.Enqueue(def)
				}
			}()
			return
		}

		return
	}

	// Success.
	record.markCompleted("Completed successfully")
	w.store.Update(record)
	w.bus.Publish(events.TaskCompleted, events.TaskCompletedEvent{
		TaskID: def.ID,
		Name:   def.Name,
	})
	w.log.Info().Str("task", def.ID).Msg("task completed")
}

// cancelTask cancels a running task by ID.
func (w *worker) cancelTask(id string) bool {
	w.mu.Lock()
	fn, ok := w.cancel[id]
	w.mu.Unlock()
	if ok {
		fn()
	}
	return ok
}
