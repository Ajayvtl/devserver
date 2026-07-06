package wsl

import (
	"context"
	"testing"
	"time"

	"github.com/Ajayvtl/devserver/internal/domain/action"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
	"github.com/Ajayvtl/devserver/internal/executor/contracts"
)

func TestWSLExecute(t *testing.T) {
	adapter := New("wsl-1", "development", "Ubuntu")

	err := adapter.Validate(context.Background(), nil)
	if err != contracts.ErrValidationFailed {
		t.Errorf("expected ErrValidationFailed for nil action, got %v", err)
	}

	targetRef, _ := valueobjects.NewReference("ls")
	act := &action.Action{
		Target:    targetRef,
		Arguments: []string{},
	}
	err = adapter.Validate(context.Background(), act)
	if err != nil {
		t.Errorf("expected no error for valid action, got %v", err)
	}
}

func TestWSLExecuteTimeout(t *testing.T) {
	adapter := New("wsl-test", "development", "Ubuntu")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	ok, diag, _ := adapter.Detect(ctx)
	if !ok && len(diag) == 0 {
		t.Errorf("expected diagnostics if detection fails")
	}
}
