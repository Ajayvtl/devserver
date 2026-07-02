package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Config struct {
	Service string
	Output  io.Writer
	Level   string
}

func New(cfg Config) zerolog.Logger {
	output := cfg.Output
	if output == nil {
		output = os.Stdout
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano

	builder := zerolog.New(output).With().Timestamp()
	if cfg.Service != "" {
		builder = builder.Str("service", cfg.Service)
	}

	log := builder.Logger()
	if cfg.Level != "" {
		if level, err := zerolog.ParseLevel(cfg.Level); err == nil {
			zerolog.SetGlobalLevel(level)
		}
	}

	return log
}

func WithContext(ctx context.Context, log zerolog.Logger) context.Context {
	return log.WithContext(ctx)
}

func FromContext(ctx context.Context) *zerolog.Logger {
	return zerolog.Ctx(ctx)
}
