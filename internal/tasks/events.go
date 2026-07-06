package tasks

import "github.com/Ajayvtl/devserver/internal/events"

// publishTaskQueued emits the queued lifecycle event.
func publishTaskQueued(bus events.Bus, task *Task) {
	if bus == nil || task == nil {
		return
	}
	bus.Publish(events.TaskQueued, events.TaskQueuedEvent{
		TaskID:      task.ID,
		Name:        task.Name,
		WorkspaceID: task.WorkspaceID,
		Detail:      "Queued",
	})
}

func publishTaskCancelled(bus events.Bus, task *Task) {
	if bus == nil || task == nil {
		return
	}
	bus.Publish(events.TaskCancelled, events.TaskCancelledEvent{
		TaskID: task.ID,
		Name:   task.Name,
	})
}

func publishTaskRolledBack(bus events.Bus, task *Task) {
	if bus == nil || task == nil {
		return
	}
	bus.Publish(events.TaskRolledBack, events.TaskRolledBackEvent{
		TaskID: task.ID,
		Name:   task.Name,
	})
}
