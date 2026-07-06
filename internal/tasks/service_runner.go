package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/Ajayvtl/devserver/internal/providers"
)

// ServiceRunner handles start, stop, restart, install, etc. for providers.
type ServiceRunner struct {
	Manager *providers.Manager
}

func (r *ServiceRunner) Execute(ctx context.Context, task *Task, runtime *Runtime) error {
	action := task.Type // e.g. "service.restart"

	target, ok := task.Payload["target"].(string)
	if !ok || target == "" {
		return fmt.Errorf("target service not specified in payload")
	}

	runtime.Logger.Info().Str("task_id", task.ID).Str("action", action).Str("target", target).Msg("Starting service runner")

	if r.Manager == nil {
		return fmt.Errorf("provider manager not configured")
	}

	provider, err := r.Manager.Get(target)
	if err != nil {
		return err
	}

	switch action {
	case "service.start":
		return provider.Start(ctx)
	case "service.stop":
		return provider.Stop(ctx)
	case "service.restart":
		// Simulating progress
		task.Progress = 25
		time.Sleep(500 * time.Millisecond) // Simulated time
		if err := provider.Stop(ctx); err != nil {
			runtime.Logger.Warn().Err(err).Msg("Stop failed during restart, continuing to start")
		}
		task.Progress = 75
		time.Sleep(500 * time.Millisecond) // Simulated time
		return provider.Start(ctx)
	case "service.install":
		return provider.Install(ctx)
	case "service.update":
		return provider.Update(ctx)
	case "service.configure":
		config, _ := task.Payload["config"].(map[string]any)
		return provider.Configure(ctx, config)
	case "service.remove":
		return provider.Uninstall(ctx)
	default:
		return fmt.Errorf("unsupported service action: %s", action)
	}
}
