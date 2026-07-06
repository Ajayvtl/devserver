package providerconfig

import "time"

// ProviderType identifies the underlying external provider.
type ProviderType string

const (
	TypeOpenAI      ProviderType = "openai"
	TypeGemini      ProviderType = "gemini"
	TypeAnthropic   ProviderType = "anthropic"
	TypeGroq        ProviderType = "groq"
	TypeOllama      ProviderType = "ollama"
	TypeAzureOpenAI ProviderType = "azure_openai"
	TypeGitHub      ProviderType = "github"
)

// ProviderConfig represents a user or organization's configured connection to an external provider.
type ProviderConfig struct {
	ID           string       `json:"id"`
	Scope        string       `json:"scope"` // "user" or "org"
	OwnerID      string       `json:"ownerId"`
	Type         ProviderType `json:"type"`
	Name         string       `json:"name"`
	Enabled      bool         `json:"enabled"`
	BaseURL      string       `json:"baseUrl"`      // For self-hosted/local overrides or custom enterprise endpoints
	DefaultModel string       `json:"defaultModel"` // e.g., "gpt-4", "claude-3-opus"
	SecretID     string       `json:"secretId"`     // Reference to the environments.Secret (API Key), never raw value
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}
