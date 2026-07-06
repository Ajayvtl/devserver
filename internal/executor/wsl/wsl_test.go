package wsl

import (
	"context"
	"testing"

	"github.com/Ajayvtl/devserver/internal/domain/action"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
	"github.com/Ajayvtl/devserver/internal/executor/contracts"
)

func TestWSLAdapter_Validate(t *testing.T) {
	adapter := New("wsl-1", "Ubuntu")
	
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

func TestWSLAdapter_Detect(t *testing.T) {
	adapter := New("wsl-test", "Ubuntu")
	ok, diag, _ := adapter.Detect(context.Background())
	if !ok && len(diag) == 0 {
		t.Errorf("expected diagnostics if detection fails")
	}
}
