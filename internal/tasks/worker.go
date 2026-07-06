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
	id       int
	log      zerolog.Logger
	bus      events.Bus
	registry *Registry
	runtime  *Runtime
	queue    *queue
	store    *store
	cancel   map[string]context.CancelFunc
	mu       sync.Mutex
}

func newWorker(id int, log zerolog.Logger, bus events.Bus, registry *Registry, runtime *Runtime, q *queue, s *store) *worker {
	return &worker{
		id:       id,
		log:      log.With().Int("worker", id).Logger(),
		bus:      bus,
		registry: registry,
		runtime:  runtime,
		queue:    q,
		store:    s,
		cancel:   make(map[string]context.CancelFunc),
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
		task := w.queue.Dequeue()
		if task == nil {
			return
		}
		w.execute(ctx, task)
	}
}

func (w *worker) execute(parentCtx context.Context, task *Task) {
	// Look up the runner for this task type.
	runner, ok := w.registry.Runner(task.Type)
	if !ok {
		task.Status = TaskFailed
		task.Error = "no runner found for task type: " + task.Type
		w.store.Update(task)
		w.bus.Publish(events.TaskFailed, events.TaskFailedEvent{
			TaskID: task.ID,
			Name:   task.Name,
			Error:  task.Error,
		})
		w.log.Error().Str("task", task.ID).Str("type", task.Type).Msg("no runner registered")
		return
	}

	// Set up cancellable context with optional timeout.
	taskCtx, cancel := context.WithCancel(parentCtx)
	if task.Timeout > 0 {
		taskCtx, cancel = context.WithTimeout(parentCtx, task.Timeout)
	}
	defer cancel()

	w.mu.Lock()
	w.cancel[task.ID] = cancel
	w.mu.Unlock()
	defer func() {
		w.mu.Lock()
		delete(w.cancel, task.ID)
		w.mu.Unlock()
	}()

	// Transition to running.
	now := time.Now().UTC()
	task.Status = TaskRunning
	task.StartedAt = now
	task.RetryCount++
	w.store.Update(task)
	w.bus.Publish(events.TaskStarted, events.TaskStartedEvent{
		TaskID: task.ID,
		Name:   task.Name,
	})
	w.log.Info().Str("task", task.ID).Str("name", task.Name).Msg("task started")

	// Execute the work via the specific runner.
	err := runner.Execute(taskCtx, task, w.runtime)

	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			task.Status = TaskCancelled
			w.store.Update(task)
			publishTaskCancelled(w.bus, task)
			w.log.Warn().Err(err).Str("task", task.ID).Msg("task cancelled")
			return
		}

		// Handle retry logic if applicable.
		if task.RetryCount <= task.MaxRetries && task.MaxRetries > 0 {
			backoff := time.Second * time.Duration(1<<task.RetryCount) // exponential backoff
			task.Status = TaskQueued
			w.store.Update(task)
			w.log.Info().Str("task", task.ID).Dur("backoff", backoff).Msg("scheduling retry")

			go func() {
				select {
				case <-parentCtx.Done():
					return
				case <-time.After(backoff):
					w.queue.Enqueue(task)
				}
			}()
			return
		}

		// Fail the task.
		task.Status = TaskFailed
		task.Error = err.Error()
		w.bus.Publish(events.TaskFailed, events.TaskFailedEvent{
			TaskID: task.ID,
			Name:   task.Name,
			Error:  task.Error,
		})
		w.store.Update(task)
		w.log.Error().Err(err).Str("task", task.ID).Msg("task failed")
		return
	}

	// Success.
	tNow := time.Now().UTC()
	task.FinishedAt = &tNow
	task.Status = TaskCompleted
	w.store.Update(task)
	w.bus.Publish(events.TaskCompleted, events.TaskCompletedEvent{
		TaskID: task.ID,
		Name:   task.Name,
	})
	w.log.Info().Str("task", task.ID).Msg("task completed")
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
