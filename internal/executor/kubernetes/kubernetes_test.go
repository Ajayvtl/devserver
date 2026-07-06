package kubernetes

import (
	"context"
	"testing"

	"github.com/Ajayvtl/devserver/internal/domain/action"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
	"github.com/Ajayvtl/devserver/internal/executor/contracts"
)

func TestKubernetesAdapter_Validate(t *testing.T) {
	adapter := New("k8s-1", "development", "default")

	err := adapter.Validate(context.Background(), nil)
	if err != contracts.ErrValidationFailed {
		t.Errorf("expected ErrValidationFailed for nil action, got %v", err)
	}

	targetRef, _ := valueobjects.NewReference("ls")
	act := &action.Action{
		Target:    targetRef,
		Arguments: []string{}, // missing pod name
	}
	err = adapter.Validate(context.Background(), act)
	if err != contracts.ErrValidationFailed {
		t.Errorf("expected ErrValidationFailed for missing pod name, got %v", err)
	}

	act.Arguments = []string{"my-pod"}
	err = adapter.Validate(context.Background(), act)
	if err != nil {
		t.Errorf("expected no error for valid action, got %v", err)
	}
}

func TestKubernetesAdapter_Detect(t *testing.T) {
	adapter := New("k8s-test", "development", "default")
	ok, diag, _ := adapter.Detect(context.Background())
	if !ok && len(diag) == 0 {
		t.Errorf("expected diagnostics if detection fails")
	}
}
