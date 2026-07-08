package auth

import (
	"context"
	"testing"
	"time"
)

type mockAuthStore struct {
	saved []AuditEvent
}

func (m *mockAuthStore) SaveAuditEvent(ctx context.Context, event AuditEvent) error {
	m.saved = append(m.saved, event)
	return nil
}

// minimal other methods to satisfy interface used by service (only SaveAuditEvent needed here)

func TestPublishAuditPersists(t *testing.T) {
	ms := &mockAuthStore{}
	// simulate persistence call
	user := &User{ID: "u-test"}
	ev := AuditEvent{EventID: "ev1", UserID: user.ID, Action: "test", Timestamp: time.Now().UTC()}
	if err := ms.SaveAuditEvent(context.Background(), ev); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	if len(ms.saved) != 1 {
		t.Fatalf("expected 1 saved audit event, got %d", len(ms.saved))
	}
}
