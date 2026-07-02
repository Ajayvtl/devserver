package executor

import (
	"context"
	"time"

	"github.com/Ajayvtl/devserver/internal/logger"
	"github.com/rs/zerolog"
)

type Command struct {
	Name    string
	Args    []string
	Dir     string
	Env     []string
	Sudo    bool
	DryRun  bool
	Timeout time.Duration
	Host    string
}

type Result struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

type Runner interface {
	Run(context.Context, Command) (Result, error)
}

type NoopRunner struct {
	log zerolog.Logger
}

func NewNoop(log zerolog.Logger) *NoopRunner {
	return &NoopRunner{log: log}
}

func (r *NoopRunner) Run(ctx context.Context, cmd Command) (Result, error) {
	log := r.log
	if ctxLog := logger.FromContext(ctx); ctxLog != nil {
		log = *ctxLog
	}

	log.Info().
		Str("command", cmd.Name).
		Strs("args", cmd.Args).
		Str("dir", cmd.Dir).
		Str("host", cmd.Host).
		Bool("sudo", cmd.Sudo).
		Bool("dry_run", cmd.DryRun).
		Msg("executor noop")

	return Result{}, nil
}
