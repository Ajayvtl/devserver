package core

import (
	"context"

	"github.com/Ajayvtl/devserver/internal/logger"
	"github.com/rs/zerolog"
)

type Task interface {
	Name() string
	Run(context.Context) error
	Rollback(context.Context) error
}

type Runner struct {
	log zerolog.Logger
}

func NewRunner(log zerolog.Logger) *Runner {
	return &Runner{log: log}
}

func (r *Runner) Run(ctx context.Context, tasks ...Task) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	log := logger.FromContext(ctx)
	completed := make([]Task, 0, len(tasks))

	for _, task := range tasks {
		log.Info().Str("task", task.Name()).Msg("task started")

		if err := task.Run(ctx); err != nil {
			log.Error().Err(err).Str("task", task.Name()).Msg("task failed")
			r.rollback(ctx, completed)
			return err
		}

		completed = append(completed, task)
		log.Info().Str("task", task.Name()).Msg("task completed")
	}

	return nil
}

func (r *Runner) rollback(ctx context.Context, tasks []Task) {
	log := logger.FromContext(ctx)
	for i := len(tasks) - 1; i >= 0; i-- {
		task := tasks[i]
		if err := task.Rollback(ctx); err != nil {
			log.Error().Err(err).Str("task", task.Name()).Msg("task rollback failed")
			continue
		}

		log.Info().Str("task", task.Name()).Msg("task rolled back")
	}
}
