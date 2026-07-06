package environments

import "time"

// EnvType categorizes the environment context.
type EnvType string

const (
	EnvDevelopment EnvType = "development"
	EnvStaging     EnvType = "staging"
	EnvProduction  EnvType = "production"
	EnvCustom      EnvType = "custom"
)

// Environment encapsulates isolated configuration sets and secrets.
type Environment struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"ownerId"` // E.g., ProjectID or OrgID
	Name      string    `json:"name"`
	Type      EnvType   `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Variable represents a plaintext environment variable.
type Variable struct {
	ID        string    `json:"id"`
	EnvID     string    `json:"envId"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Secret represents an encrypted credential (e.g., API keys, Tokens).
// The Value field is NEVER serialized to plaintext automatically.
type Secret struct {
	ID        string    `json:"id"`
	EnvID     string    `json:"envId"`
	Key       string    `json:"key"`
	Value     string    `json:"-"` // Explicitly omitted from JSON to prevent accidental exposure
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SecretResponse is used when a secure abstraction reference is returned.
type SecretReference struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}
