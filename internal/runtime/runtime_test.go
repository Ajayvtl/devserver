package runtime_test

import (
	"context"
	"testing"
	"time"

	"github.com/Ajayvtl/devserver/internal/runtime"
)

// MockComponent simulates an existing subsystem (e.g., Task Engine)
// adapting to the Component interface.
type MockComponent struct {
	name   string
	status runtime.Status
	health runtime.Health
}

func NewMockComponent(name string) *MockComponent {
	return &MockComponent{
		name:   name,
		status: runtime.StatusStopped,
		health: runtime.HealthHealthy,
	}
}

func (m *MockComponent) Name() string {
	return m.name
}

func (m *MockComponent) Initialize(ctx context.Context) error {
	m.status = runtime.StatusStarting
	return nil
}

func (m *MockComponent) Start(ctx context.Context) error {
	m.status = runtime.StatusRunning
	return nil
}

func (m *MockComponent) Stop(ctx context.Context) error {
	m.status = runtime.StatusStopped
	return nil
}

func (m *MockComponent) Status() runtime.Status {
	return m.status
}

func (m *MockComponent) Health() runtime.Health {
	return m.health
}

func TestRuntimeBootstrap(t *testing.T) {
	registry := runtime.NewRegistry()

	// 1. Register components (simulating existing subsystems adapting to Component)
	tasksEngine := NewMockComponent("tasks.Engine")
	providerManager := NewMockComponent("providers.Manager")

	if err := registry.Register(tasksEngine); err != nil {
		t.Fatalf("Failed to register tasksEngine: %v", err)
	}
	if err := registry.Register(providerManager); err != nil {
		t.Fatalf("Failed to register providerManager: %v", err)
	}

	coordinator := runtime.NewCoordinator(registry)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 2. Initialize
	if err := coordinator.Initialize(ctx); err != nil {
		t.Fatalf("Coordinator Initialize failed: %v", err)
	}
	if tasksEngine.Status() != runtime.StatusStarting {
		t.Errorf("Expected tasksEngine status to be Starting, got %v", tasksEngine.Status())
	}

	// 3. Start
	if err := coordinator.Start(ctx); err != nil {
		t.Fatalf("Coordinator Start failed: %v", err)
	}
	if tasksEngine.Status() != runtime.StatusRunning {
		t.Errorf("Expected tasksEngine status to be Running, got %v", tasksEngine.Status())
	}

	// 4. Health
	if tasksEngine.Health() != runtime.HealthHealthy {
		t.Errorf("Expected tasksEngine health to be Healthy, got %v", tasksEngine.Health())
	}

	// 5. Stop
	if err := coordinator.Stop(ctx); err != nil {
		t.Fatalf("Coordinator Stop failed: %v", err)
	}
	if tasksEngine.Status() != runtime.StatusStopped {
		t.Errorf("Expected tasksEngine status to be Stopped, got %v", tasksEngine.Status())
	}
}
