package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/Ajayvtl/devserver/internal/bootstrap"
	"github.com/Ajayvtl/devserver/internal/config"
	"github.com/Ajayvtl/devserver/internal/core"
	"github.com/Ajayvtl/devserver/internal/executor"
	"github.com/Ajayvtl/devserver/internal/logger"
	"github.com/Ajayvtl/devserver/internal/modules/mysql"
	"github.com/Ajayvtl/devserver/internal/modules/nginx"
	"github.com/Ajayvtl/devserver/internal/modules/node"
	"github.com/Ajayvtl/devserver/internal/modules/php"
	"github.com/Ajayvtl/devserver/internal/modules/postgres"
	"github.com/Ajayvtl/devserver/internal/modules/python"
	"github.com/Ajayvtl/devserver/internal/modules/redis"
	"github.com/Ajayvtl/devserver/internal/platform"
	"github.com/Ajayvtl/devserver/internal/registry"
	"github.com/Ajayvtl/devserver/internal/state"
	"github.com/rs/zerolog"
)

func Run(ctx context.Context, args []string) error {
	log := logger.New(logger.Config{
		Service: "devserver",
		Output:  os.Stdout,
	})

	ctx = logger.WithContext(ctx, log)

	if len(args) > 0 && args[0] == "serve" {
		return runServer(ctx, log)
	}

	cfg, err := config.NewLoader().Load(ctx, "")
	if err != nil {
		return err
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
	exec := executor.NewNoop(log)
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

func runServer(ctx context.Context, log zerolog.Logger) error {
	db, err := state.NewDB("configs/devserver.db")
	if err != nil {
		return err
	}
	defer db.Close()

	root, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	indexer := core.NewIndexer(log, []core.WorkspaceSpec{{ID: "devserver", Root: root}})
	server := core.NewAPIServer(log, db, indexer, ":8080")
	return server.Run(ctx)
}
