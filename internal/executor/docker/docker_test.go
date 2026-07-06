package docker

import (
	"context"
	"testing"

	"github.com/Ajayvtl/devserver/internal/domain/action"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
	"github.com/Ajayvtl/devserver/internal/executor/contracts"
)

func TestDockerAdapter_Validate(t *testing.T) {
	adapter := New("docker-1", "development")

	err := adapter.Validate(context.Background(), nil)
	if err != contracts.ErrValidationFailed {
		t.Errorf("expected ErrValidationFailed for nil action, got %v", err)
	}

	targetRef, _ := valueobjects.NewReference("ls")
	act := &action.Action{
		Target:    targetRef,
		Arguments: []string{}, // No container id
	}
	err = adapter.Validate(context.Background(), act)
	if err != contracts.ErrValidationFailed {
		t.Errorf("expected ErrValidationFailed for missing container id, got %v", err)
	}

	act.Arguments = []string{"my-container"}
	err = adapter.Validate(context.Background(), act)
	if err != nil {
		t.Errorf("expected no error for valid action, got %v", err)
	}
}

func TestDockerAdapter_Detect(t *testing.T) {
	// We only verify it doesn't panic and returns a valid diagnostic since docker might not be running
	adapter := New("docker-test", "development")
	ok, diag, _ := adapter.Detect(context.Background())
	if !ok && len(diag) == 0 {
		t.Errorf("expected diagnostics if detection fails")
	}
}
