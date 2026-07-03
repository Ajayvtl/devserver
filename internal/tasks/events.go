package tasks

import "github.com/Ajayvtl/devserver/internal/events"

// publishTaskQueued emits the queued lifecycle event.
func publishTaskQueued(bus events.Bus, record *Record) {
	if bus == nil || record == nil {
		return
	}
	bus.Publish(events.TaskQueued, events.TaskQueuedEvent{
		TaskID:      record.ID,
		Name:        record.Name,
		WorkspaceID: record.WorkspaceID,
		Detail:      record.Detail,
	})
}

func publishTaskCancelled(bus events.Bus, record *Record) {
	if bus == nil || record == nil {
		return
	}
	bus.Publish(events.TaskCancelled, events.TaskCancelledEvent{
		TaskID: record.ID,
		Name:   record.Name,
	})
}

func publishTaskRolledBack(bus events.Bus, record *Record) {
	if bus == nil || record == nil {
		return
	}
	bus.Publish(events.TaskRolledBack, events.TaskRolledBackEvent{
		TaskID: record.ID,
		Name:   record.Name,
	})
}
