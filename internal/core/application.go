package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/Ajayvtl/devserver/internal/config"
	"github.com/Ajayvtl/devserver/internal/executor/legacy"
	"github.com/Ajayvtl/devserver/internal/logger"
	"github.com/Ajayvtl/devserver/internal/platform"
	"github.com/Ajayvtl/devserver/internal/state"
	"github.com/rs/zerolog"
)

type Bootstrapper interface {
	Bootstrap(context.Context) error
}

type Application struct {
	log       zerolog.Logger
	registry  *Registry
	bootstrap Bootstrapper
	cfg       config.Config
	store     state.Store
	detector  platform.Detector
	executor  legacy.Runner
	lifecycle *Lifecycle
}

type Dependencies struct {
	Logger    zerolog.Logger
	Registry  *Registry
	Bootstrap Bootstrapper
	Config    config.Config
	State     state.Store
	Detector  platform.Detector
	Executor  legacy.Runner
}

func New(deps Dependencies) *Application {
	runner := NewRunner(deps.Logger)

	return &Application{
		log:       deps.Logger,
		registry:  deps.Registry,
		bootstrap: deps.Bootstrap,
		cfg:       deps.Config,
		store:     deps.State,
		detector:  deps.Detector,
		executor:  deps.Executor,
		lifecycle: NewLifecycle(runner),
	}
}

func (a *Application) Run(ctx context.Context, args []string) error {
	ctx = logger.WithContext(ctx, a.log)
	log := logger.FromContext(ctx)

	log.Info().Strs("args", args).Msg("command received")

	if len(args) == 0 {
		return a.handleDefault(ctx)
	}

	switch args[0] {
	case "doctor":
		return a.doctor(ctx)
	case "install":
		return a.bootstrap.Bootstrap(ctx)
	case "service":
		return a.service(ctx, args[1:])
	case "deploy":
		return a.deploy(ctx, args[1:])
	case "backup", "config", "monitor":
		log.Info().Str("command", args[0]).Msg("command scaffolded")
		return nil
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func (a *Application) handleDefault(ctx context.Context) error {
	log := logger.FromContext(ctx)
	if a.cfg.Commands.Default == "" {
		log.Info().Msg("no default command configured")
		return nil
	}

	return a.Run(ctx, []string{a.cfg.Commands.Default})
}

func (a *Application) doctor(ctx context.Context) error {
	log := logger.FromContext(ctx)

	if a.detector != nil {
		platformInfo, err := a.detector.Detect(ctx)
		if err != nil {
			return err
		}

		log.Info().
			Str("platform_id", platformInfo.ID).
			Str("platform_name", platformInfo.Name).
			Str("platform_version", platformInfo.Version).
			Msg("platform detected")
	}

	if a.store != nil {
		st, err := a.store.Load(ctx)
		if err != nil {
			return err
		}

		log.Info().Int("installed_modules", len(st.Installed)).Msg("state loaded")
	}

	return nil
}

func (a *Application) service(ctx context.Context, args []string) error {
	log := logger.FromContext(ctx)
	if len(args) < 2 {
		log.Info().Msg("service command scaffolded")
		return nil
	}

	action := args[0]
	moduleName := args[1]
	if a.registry == nil {
		return fmt.Errorf("module registry is not configured")
	}

	module, ok := a.registry.Lookup(moduleName)
	if !ok {
		return fmt.Errorf("module %q is not registered", moduleName)
	}

	log.Info().
		Str("action", action).
		Str("module", module.Name()).
		Msg("service command routed")
	return nil
}

func (a *Application) deploy(ctx context.Context, args []string) error {
	log := logger.FromContext(ctx)
	log.Info().Strs("args", args).Msg("deploy command scaffolded")
	return nil
}

func (a *Application) RegisteredModules() []string {
	modules := a.registry.Modules()
	out := make([]string, 0, len(modules))
	for _, module := range modules {
		out = append(out, module.Name())
	}
	return out
}

func (a *Application) ServiceName() string {
	return strings.TrimSpace(a.cfg.Service.Name)
}
