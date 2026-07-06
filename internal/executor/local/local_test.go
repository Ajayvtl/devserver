package local

import (
	"context"
	"testing"
	"time"

	"github.com/Ajayvtl/devserver/internal/domain/action"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
)

func TestLocalExecute(t *testing.T) {
	adapter := New("test-local-1", "development")

	targetRef, _ := valueobjects.NewReference("hostname")
	act := &action.Action{
		Target:    targetRef,
		Arguments: []string{},
	}

	res, err := adapter.Execute(context.Background(), act)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", res.ExitCode)
	}
}

func TestLocalExecuteTimeout(t *testing.T) {
	adapter := New("test-local-timeout", "development")

	targetRef, _ := valueobjects.NewReference("powershell")
	act := &action.Action{
		Target:    targetRef,
		Arguments: []string{"-Command", "Start-Sleep -Seconds 2"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := adapter.Execute(ctx, act)
	if err == nil {
		t.Fatalf("expected timeout error, got nil")
	}
}
