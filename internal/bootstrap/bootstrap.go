package bootstrap

import (
	"context"
	"fmt"

	"github.com/Ajayvtl/devserver/internal/executor/legacy"
	"github.com/Ajayvtl/devserver/internal/logger"
	"github.com/Ajayvtl/devserver/internal/platform"
	"github.com/Ajayvtl/devserver/internal/registry"
	"github.com/Ajayvtl/devserver/internal/state"
	"github.com/rs/zerolog"
)

type Bootstrapper struct {
	log      zerolog.Logger
	registry *registry.Registry
	store    state.Store
	exec     legacy.Runner
	platform platform.Detector
}

func New(log zerolog.Logger, reg *registry.Registry, store state.Store, exec legacy.Runner, detector platform.Detector) *Bootstrapper {
	return &Bootstrapper{
		log:      log,
		registry: reg,
		store:    store,
		exec:     exec,
		platform: detector,
	}
}

func (b *Bootstrapper) Bootstrap(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	log := b.log
	if ctxLog := logger.FromContext(ctx); ctxLog != nil {
		log = *ctxLog
	}

	log.Info().Msg("bootstrap started")

	if b.platform != nil {
		platformInfo, err := b.platform.Detect(ctx)
		if err != nil {
			return err
		}
		log.Info().Str("platform", platformInfo.ID).Msg("platform detected")
	}

	current := state.New()
	if b.store != nil {
		st, err := b.store.Load(ctx)
		if err != nil {
			return err
		}
		current = st
	}

	if b.registry == nil {
		log.Info().Msg("no modules registered")
		return nil
	}

	for _, module := range b.registry.Modules() {
		if err := module.Check(ctx); err != nil {
			return fmt.Errorf("%s check failed: %w", module.Name(), err)
		}
		if err := module.Validate(ctx); err != nil {
			return fmt.Errorf("%s validate failed: %w", module.Name(), err)
		}
		if err := module.Install(ctx); err != nil {
			return fmt.Errorf("%s install failed: %w", module.Name(), err)
		}
		if err := module.Configure(ctx); err != nil {
			return fmt.Errorf("%s configure failed: %w", module.Name(), err)
		}

		current.MarkInstalled(module.Name(), "unknown")
	}

	if b.store != nil {
		if err := b.store.Save(ctx, current); err != nil {
			return err
		}
	}

	log.Info().Msg("bootstrap finished")
	return nil
}
