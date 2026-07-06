package settings

import "time"

// Scope defines whether a setting or integration applies to a user or an organization.
type Scope string

const (
	ScopeUser Scope = "user"
	ScopeOrg  Scope = "org"
)

// Setting represents a key-value configuration pair for a specific scope.
type Setting struct {
	ID        string    `json:"id"`
	Scope     Scope     `json:"scope"`
	OwnerID   string    `json:"ownerId"` // UserID or OrgID based on Scope
	Key       string    `json:"key"`
	Value     string    `json:"value"` // Stored as a JSON string for flexibility
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Integration represents a configured third-party integration (provider-agnostic).
type Integration struct {
	ID        string    `json:"id"`
	Scope     Scope     `json:"scope"`
	OwnerID   string    `json:"ownerId"`
	Provider  string    `json:"provider"` // e.g., "github", "slack", "jira"
	Status    string    `json:"status"`   // e.g., "active", "inactive", "pending"
	Config    string    `json:"config"`   // JSON-encoded metadata/configuration
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
