package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/Ajayvtl/devserver/internal/ai"
	"github.com/Ajayvtl/devserver/internal/api"
	"github.com/Ajayvtl/devserver/internal/auth"
	"github.com/Ajayvtl/devserver/internal/bootstrap"
	"github.com/Ajayvtl/devserver/internal/capabilities"
	"github.com/Ajayvtl/devserver/internal/commands"
	"github.com/Ajayvtl/devserver/internal/config"
	"github.com/Ajayvtl/devserver/internal/core"
	"github.com/Ajayvtl/devserver/internal/environments"
	"github.com/Ajayvtl/devserver/internal/events"
	"github.com/Ajayvtl/devserver/internal/executor/legacy"
	"github.com/Ajayvtl/devserver/internal/filesystem"
	"github.com/Ajayvtl/devserver/internal/logger"
	"github.com/Ajayvtl/devserver/internal/modules/mysql"
	"github.com/Ajayvtl/devserver/internal/modules/nginx"
	"github.com/Ajayvtl/devserver/internal/modules/node"
	"github.com/Ajayvtl/devserver/internal/modules/php"
	"github.com/Ajayvtl/devserver/internal/modules/postgres"
	"github.com/Ajayvtl/devserver/internal/modules/python"
	"github.com/Ajayvtl/devserver/internal/modules/redis"
	"github.com/Ajayvtl/devserver/internal/platform"
	"github.com/Ajayvtl/devserver/internal/providerconfig"
	"github.com/Ajayvtl/devserver/internal/providers"
	"github.com/Ajayvtl/devserver/internal/rbac"
	"github.com/Ajayvtl/devserver/internal/registry"
	rt "github.com/Ajayvtl/devserver/internal/runtime"
	"github.com/Ajayvtl/devserver/internal/settings"
	"github.com/Ajayvtl/devserver/internal/state"
	"github.com/Ajayvtl/devserver/internal/tasks"
	"github.com/rs/zerolog"
)

func Run(ctx context.Context, args []string) error {
	log := logger.New(logger.Config{
		Service: "devserver",
		Output:  os.Stdout,
	})

	ctx = logger.WithContext(ctx, log)

	cfg, err := config.NewLoader().Load(ctx, "")
	if err != nil {
		return err
	}

	if len(args) > 0 && args[0] == "serve" {
		return runServer(ctx, log, cfg)
	}

	reg := core.NewRegistry()
	for _, module := range []registry.Module{
		&nginx.Module{},
		&node.Module{},
		&postgres.Module{},
		&mysql.Module{},
		&redis.Module{},
		&php.Module{},
		&python.Module{},
	} {
		if err := reg.Register(module); err != nil {
			return err
		}
	}

	store := state.NewMemoryStore()
	detector := platform.NewDetector()
	exec := legacy.NewNoop(log)
	boot := bootstrap.New(log, reg.Inner(), store, exec, detector)

	application := core.New(core.Dependencies{
		Logger:    log,
		Registry:  reg,
		Bootstrap: boot,
		Config:    cfg,
		State:     store,
		Detector:  detector,
		Executor:  exec,
	})

	return application.Run(ctx, args)
}

func runServer(ctx context.Context, log zerolog.Logger, cfg config.Config) error {
	db, err := state.NewDB("configs/devserver.db")
	if err != nil {
		return err
	}
	defer db.Close()

	root, err := filepath.Abs(".")
	if err != nil {
		return err
	}

	// Core runtime infrastructure.
	bus := events.NewBus()
	localFS := filesystem.NewLocal()
	taskStream := tasks.NewEventStream(log, bus)

	// Create subsystems — they start empty, Manager fills them.
	indexer := core.NewIndexer(log, nil, bus)
	providerManager := providers.NewManager()
	provider := core.NewWorkspaceProvider(indexer, providerManager)
	capRegistry := capabilities.NewRegistry()

	// Initialize Provider Runtime (Milestone 3)
	execRuntime := legacy.NewLocalRuntime()
	providerManager.Register(providers.NewRedisProvider(execRuntime))
	providerManager.Register(providers.NewNodeProvider(execRuntime))
	providerManager.Register(providers.NewAIProvider(execRuntime))
	providerManager.Register(&providers.LocalProvider{
		Runtime:   execRuntime,
		Meta:      providers.ProviderMetadata{Name: "docker", Description: "Docker container runtime", Version: "latest"},
		CheckCmd:  "docker",
		CheckArgs: []string{"--version"},
	})
	providerManager.Register(&providers.LocalProvider{
		Runtime:   execRuntime,
		Meta:      providers.ProviderMetadata{Name: "mysql", Description: "MySQL relational database", Version: "latest"},
		CheckCmd:  "mysql",
		CheckArgs: []string{"--version"},
	})
	providerManager.Register(&providers.LocalProvider{
		Runtime:   execRuntime,
		Meta:      providers.ProviderMetadata{Name: "nginx", Description: "Nginx web server", Version: "latest"},
		CheckCmd:  "nginx",
		CheckArgs: []string{"-v"},
	})

	// Initialize Task Engine
	taskRegistry := tasks.NewRegistry()
	taskRegistry.Register("workspace.index", &tasks.WorkspaceIndexRunner{})
	taskRegistry.Register("workspace.setup", &tasks.WorkspaceSetupRunner{})

	serviceRunner := &tasks.ServiceRunner{Manager: providerManager}
	taskRegistry.Register("service.start", serviceRunner)
	taskRegistry.Register("service.stop", serviceRunner)
	taskRegistry.Register("service.restart", serviceRunner)
	taskRegistry.Register("service.install", serviceRunner)
	taskRegistry.Register("service.update", serviceRunner)
	taskRegistry.Register("service.configure", serviceRunner)
	taskRegistry.Register("service.remove", serviceRunner)

	// Phase 6: AI Runtime
	aiRuntime := ai.NewRuntime(log, provider, indexer, providerManager)

	// AI Task Runners
	aiRunner := &ai.TaskRunner{AIRuntime: aiRuntime, Bus: bus}
	taskRegistry.Register("ai.generate", aiRunner)
	taskRegistry.Register("ai.refactor", aiRunner)
	taskRegistry.Register("ai.explain", aiRunner)

	runtime := &tasks.Runtime{
		Logger:   log,
		Provider: provider,
		Indexer:  indexer,
		EventBus: taskStream,
		Registry: capRegistry,
		// Store and Commands to be set
	}
	taskEngine := tasks.NewEngine(log, bus, taskRegistry, runtime, tasks.EngineConfig{Workers: 4, QueueMax: 128})

	var watcher *core.WorkspaceWatcher
	watcher, err = core.NewWorkspaceWatcher(log, bus)
	if err != nil {
		log.Warn().Err(err).Msg("filesystem watcher unavailable, continuing without live watch")
	}

	modules := []registry.Module{
		&nginx.Module{},
		&node.Module{},
		&postgres.Module{},
		&mysql.Module{},
		&redis.Module{},
		&php.Module{},
		&python.Module{},
	}
	for _, module := range modules {
		if err := capRegistry.Register(capabilities.NewModuleProvider(module)); err != nil {
			return err
		}
	}

	commandEngine := commands.NewEngine(log, bus, capRegistry, taskEngine, db)

	// WorkspaceManager owns the lifecycle.
	manager := core.NewWorkspaceManager(core.WorkspaceManagerDeps{
		Logger:   log,
		Indexer:  indexer,
		Watcher:  watcher,
		Provider: provider,
		Bus:      bus,
		FS:       localFS,
	})

	// Open the current project as the default workspace.
	if err := manager.Open(ctx, "devserver", root); err != nil {
		return err
	}

	// Phase 7: Authentication
	authStore, err := auth.NewMySQLStore(cfg.Database.DSN)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize auth database")
	}
	defer authStore.Close()

	authService := auth.NewService(log, bus, authStore)
	authService.RegisterProvider(auth.NewLocalProvider(authStore))

	// Phase 7: RBAC
	rbacStore, err := rbac.NewMySQLStore(cfg.Database.DSN)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize rbac database")
	}
	defer rbacStore.Close()

	rbacService := rbac.NewService(log, rbacStore)
	_ = rbacService // Silencing unused warning until wired to HTTP handlers

	// Phase 7: Settings & Integrations
	settingsStore, err := settings.NewMySQLStore(cfg.Database.DSN)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize settings database")
	}
	defer settingsStore.Close()

	settingsService := settings.NewService(log, settingsStore)
	_ = settingsService // Silencing unused warning until wired to HTTP handlers

	// Phase 7: Environment Management
	envStore, err := environments.NewMySQLStore(cfg.Database.DSN)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize environments database")
	}
	defer envStore.Close()

	cryptoService, err := environments.NewAESCryptoService([]byte(cfg.Security.MasterKey))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize crypto service")
	}

	envService := environments.NewService(log, envStore, cryptoService)
	_ = envService // Silencing unused warning until wired to HTTP handlers

	// Phase 7: Provider Configuration
	providerConfigStore, err := providerconfig.NewMySQLStore(cfg.Database.DSN)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize providerconfig database")
	}
	defer providerConfigStore.Close()

	providerConfigService := providerconfig.NewService(log, providerConfigStore, envService)
	_ = providerConfigService // Silencing unused warning until wired to HTTP handlers

	// Phase 8: HTTP API & Middleware Integration
	apiRouter := api.NewRouter(
		authService,
		rbacService,
		settingsService,
		envService,
		providerConfigService,
	)

	server := core.NewAPIServer(log, db, indexer, provider, taskEngine, taskStream, commandEngine, capRegistry, apiRouter, ":8080")

	// Phase 2: Runtime Bootstrap
	rtReg := rt.NewRegistry()
	_ = rtReg.Register(taskStream)
	_ = rtReg.Register(providerManager)
	_ = rtReg.Register(provider)
	_ = rtReg.Register(indexer)
	_ = rtReg.Register(taskEngine)
	_ = rtReg.Register(commandEngine)
	_ = rtReg.Register(server)
	_ = rtReg.Register(aiRuntime)

	coordinator := rt.NewCoordinator(rtReg)

	if err := coordinator.Initialize(ctx); err != nil {
		log.Error().Err(err).Msg("runtime initialization failed")
		return err
	}

	if err := coordinator.Start(ctx); err != nil {
		log.Error().Err(err).Msg("runtime start failed")
		_ = coordinator.Stop(ctx)
		return err
	}

	go syncTaskStore(ctx, bus, db, log)

	if watcher != nil {
		go watcher.Run(ctx)
	}

	// Block until context is cancelled
	<-ctx.Done()

	if err := coordinator.Stop(context.Background()); err != nil {
		log.Error().Err(err).Msg("runtime stop failed")
	}

	return nil
}

func syncTaskStore(ctx context.Context, bus events.Bus, db *state.StoreDB, log zerolog.Logger) {
	if db == nil || bus == nil {
		return
	}

	subscriptions := []struct {
		eventType events.EventType
	}{
		{events.TaskStarted},
		{events.TaskProgress},
		{events.TaskCompleted},
		{events.TaskFailed},
		{events.TaskCancelled},
		{events.TaskRolledBack},
	}

	for _, sub := range subscriptions {
		ch := bus.Subscribe(sub.eventType)
		go func(topic events.EventType, eventsCh events.Subscriber) {
			for {
				select {
				case <-ctx.Done():
					return
				case evt, ok := <-eventsCh:
					if !ok {
						return
					}
					if err := applyTaskEvent(ctx, db, topic, evt.Payload); err != nil {
						log.Warn().Err(err).Str("topic", string(topic)).Msg("task event sync failed")
					}
				}
			}
		}(sub.eventType, ch)
	}
}

func applyTaskEvent(ctx context.Context, db *state.StoreDB, topic events.EventType, payload any) error {
	switch topic {
	case events.TaskStarted:
		if evt, ok := payload.(events.TaskStartedEvent); ok {
			return db.UpdateTask(ctx, evt.TaskID, 10, "Running", "Started")
		}
	case events.TaskProgress:
		if evt, ok := payload.(events.TaskProgressEvent); ok {
			return db.UpdateTask(ctx, evt.TaskID, evt.Progress, "Running", evt.Detail)
		}
	case events.TaskCompleted:
		if evt, ok := payload.(events.TaskCompletedEvent); ok {
			return db.UpdateTask(ctx, evt.TaskID, 100, "Completed", "Completed successfully")
		}
	case events.TaskFailed:
		if evt, ok := payload.(events.TaskFailedEvent); ok {
			return db.UpdateTask(ctx, evt.TaskID, 0, "Failed", evt.Error)
		}
	case events.TaskCancelled:
		if evt, ok := payload.(events.TaskCancelledEvent); ok {
			return db.UpdateTask(ctx, evt.TaskID, 0, "Cancelled", "Cancelled")
		}
	case events.TaskRolledBack:
		if evt, ok := payload.(events.TaskRolledBackEvent); ok {
			return db.UpdateTask(ctx, evt.TaskID, 0, "RolledBack", "Rolled back")
		}
	}
	return nil
}
