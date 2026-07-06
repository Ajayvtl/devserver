package ports

import (
	"context"
	"time"
)

// Clock abstracts system time to allow testing time-dependent logic.
type Clock interface {
	Now() time.Time
}

// Logger defines the application's logging interface.
type Logger interface {
	Info(ctx context.Context, msg string, args ...interface{})
	Error(ctx context.Context, msg string, err error, args ...interface{})
	Debug(ctx context.Context, msg string, args ...interface{})
}

// Notifier handles sending alerts or notifications to external systems.
type Notifier interface {
	Notify(ctx context.Context, topic string, payload interface{}) error
}

// SecretResolver securely retrieves runtime secrets without storing them.
type SecretResolver interface {
	Resolve(ctx context.Context, secretRef string) (string, error)
}

// MetricsCollector abstracts telemetry and metric aggregation.
type MetricsCollector interface {
	Record(ctx context.Context, metric string, value float64, tags map[string]string) error
}

// AuditWriter records significant domain events for compliance and history.
type AuditWriter interface {
	WriteAudit(ctx context.Context, eventType string, details interface{}) error
}
