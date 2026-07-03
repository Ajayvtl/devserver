package tasks

import "time"

// markRunning moves a record into the running state.
func (r *Record) markRunning() {
	now := time.Now().UTC()
	r.Status = StatusRunning
	r.StartedAt = &now
}

// markCompleted finalizes a record as completed.
func (r *Record) markCompleted(detail string) {
	now := time.Now().UTC()
	r.Status = StatusCompleted
	r.Progress = 100
	r.Detail = detail
	r.CompletedAt = &now
}

// markFailed finalizes a record as failed.
func (r *Record) markFailed(err string) {
	now := time.Now().UTC()
	r.Status = StatusFailed
	r.Error = err
	r.CompletedAt = &now
}

// markCancelled finalizes a record as cancelled.
func (r *Record) markCancelled(detail string) {
	now := time.Now().UTC()
	r.Status = StatusCancelled
	if detail != "" {
		r.Detail = detail
	}
	r.CompletedAt = &now
}

// markRolledBack finalizes a record as rolled back.
func (r *Record) markRolledBack(detail string) {
	now := time.Now().UTC()
	r.Status = StatusRolledBack
	if detail != "" {
		r.Detail = detail
	}
	r.CompletedAt = &now
}
