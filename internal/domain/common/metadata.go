package common

import "time"

// Metadata is the standard metadata struct embedded in all persistent entities.
type Metadata struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy string
	UpdatedBy string

	// Versioning
	Generation int // Increments on spec changes
	Revision   int // Increments on any change (status, etc)

	// Classifications
	Labels      map[string]string // Machine-readable (e.g., tier: database)
	Tags        []string          // Human-readable (e.g., "critical", "legacy")
	Annotations map[string]string // Arbitrary non-identifying metadata
}
