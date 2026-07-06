package local

import (
	"context"
	"testing"
	"time"

	"github.com/Ajayvtl/devserver/internal/domain/action"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
)

func TestLocalAdapter_Execute_Success(t *testing.T) {
	adapter := New("test-local-1")

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

func TestLocalAdapter_Execute_Timeout(t *testing.T) {
	adapter := New("test-local-timeout")

	targetRef, _ := valueobjects.NewReference("powershell")
	act := &action.Action{
		Target:    targetRef,
		Arguments: []string{"-Command", "Start-Sleep -Seconds 2"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := adapter.Execute(ctx, act)
	if err == nil {
		t.Fatalf("expected timeout error, got nil")
	}
}
