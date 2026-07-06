package ssh

import (
	"context"
	"testing"
	"time"

	"github.com/Ajayvtl/devserver/internal/domain/action"
	"github.com/Ajayvtl/devserver/internal/domain/valueobjects"
	"github.com/Ajayvtl/devserver/internal/executor/contracts"
	"golang.org/x/crypto/ssh"
)

func TestSSHAdapter_Validate(t *testing.T) {
	config := &ssh.ClientConfig{Timeout: 1 * time.Second}
	adapter := New("ssh-1", "localhost:22", config)
	
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

func TestSSHAdapter_Detect(t *testing.T) {
	config := &ssh.ClientConfig{Timeout: 10 * time.Millisecond}
	// connect to invalid port to ensure failure diagnostic
	adapter := New("ssh-test", "127.0.0.1:0", config)
	ok, diag, _ := adapter.Detect(context.Background())
	if ok {
		t.Errorf("expected detect to fail on invalid port")
	}
	if len(diag) == 0 {
		t.Errorf("expected diagnostics if detection fails")
	}
}
