package core

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/Ajayvtl/devserver/internal/logger"
)

type fakeTask struct {
	name     string
	runErr   error
	rollback *[]string
	events   *[]string
}

func (t fakeTask) Name() string { return t.name }

func (t fakeTask) Run(context.Context) error {
	*t.events = append(*t.events, "run:"+t.name)
	return t.runErr
}

func (t fakeTask) Rollback(context.Context) error {
	*t.events = append(*t.events, "rollback:"+t.name)
	*t.rollback = append(*t.rollback, t.name)
	return nil
}

func TestRunnerRollsBackCompletedTasksOnFailure(t *testing.T) {
	events := []string{}
	rolledBack := []string{}

	runner := NewRunner(logger.New(logger.Config{Output: io.Discard}))
	err := runner.Run(context.Background(),
		fakeTask{name: "first", events: &events, rollback: &rolledBack},
		fakeTask{name: "second", runErr: errors.New("boom"), events: &events, rollback: &rolledBack},
	)

	if err == nil {
		t.Fatal("expected an error")
	}

	wantEvents := []string{"run:first", "run:second", "rollback:first"}
	if len(events) != len(wantEvents) {
		t.Fatalf("unexpected events: %v", events)
	}
	for i, want := range wantEvents {
		if events[i] != want {
			t.Fatalf("unexpected event at %d: got %q want %q", i, events[i], want)
		}
	}

	if len(rolledBack) != 1 || rolledBack[0] != "first" {
		t.Fatalf("unexpected rollback order: %v", rolledBack)
	}
}
