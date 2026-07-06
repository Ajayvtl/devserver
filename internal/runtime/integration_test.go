package runtime_test

import (
	"context"
	"testing"
	"time"

	"github.com/Ajayvtl/devserver/internal/core"
	"github.com/Ajayvtl/devserver/internal/events"
	"github.com/Ajayvtl/devserver/internal/providers"
	"github.com/Ajayvtl/devserver/internal/runtime"
	"github.com/Ajayvtl/devserver/internal/tasks"
	"github.com/rs/zerolog"
)

func TestIntegrationBootstrap(t *testing.T) {
	log := zerolog.Nop()
	bus := events.NewBus()

	// Instantiate real components
	taskStream := tasks.NewEventStream(log, bus)
	providerManager := providers.NewManager()
	indexer := core.NewIndexer(log, nil, bus)
	workspaceProvider := core.NewWorkspaceProvider(indexer, providerManager)

	// Step 1: Registration
	rtReg := runtime.NewRegistry()
	if err := rtReg.Register(taskStream); err != nil {
		t.Fatalf("Failed to register taskStream: %v", err)
	}
	if err := rtReg.Register(providerManager); err != nil {
		t.Fatalf("Failed to register providerManager: %v", err)
	}
	if err := rtReg.Register(workspaceProvider); err != nil {
		t.Fatalf("Failed to register workspaceProvider: %v", err)
	}
	if err := rtReg.Register(indexer); err != nil {
		t.Fatalf("Failed to register indexer: %v", err)
	}

	coordinator := runtime.NewCoordinator(rtReg)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Step 2: Initialization
	if err := coordinator.Initialize(ctx); err != nil {
		t.Fatalf("Coordinator Initialize failed: %v", err)
	}

	// Step 3: Startup
	if err := coordinator.Start(ctx); err != nil {
		t.Fatalf("Coordinator Start failed: %v", err)
	}

	// Step 4: Health Aggregation
	health := coordinator.AggregateHealth()
	if health.Status != "healthy" {
		t.Errorf("Expected aggregated health to be 'healthy', got '%s'", health.Status)
	}
	if len(health.Components) != 4 {
		t.Errorf("Expected 4 component health entries, got %d", len(health.Components))
	}

	// Step 5: Shutdown
	if err := coordinator.Stop(context.Background()); err != nil {
		t.Fatalf("Coordinator Stop failed: %v", err)
	}
}
