package providerconfig

import (
	"context"
	"errors"
	"net/http"

	"github.com/Ajayvtl/devserver/internal/domain/common"
	"github.com/Ajayvtl/devserver/internal/environments"
	"github.com/rs/zerolog"
)

// Store defines persistence for provider configurations.
type Store interface {
	GetConfig(ctx context.Context, id string) (*ProviderConfig, error)
	ListConfigs(ctx context.Context, scope, ownerID string, params common.QueryParams) ([]*ProviderConfig, int, error)
	UpsertConfig(ctx context.Context, cfg *ProviderConfig) error
	DeleteConfig(ctx context.Context, id string) error
}

// Service manages provider configurations and handles runtime binding via environments.
type Service interface {
	SaveConfig(ctx context.Context, cfg *ProviderConfig) error
	GetConfig(ctx context.Context, id string) (*ProviderConfig, error)
	ListConfigs(ctx context.Context, scope, ownerID string, params common.QueryParams) ([]*ProviderConfig, int, error)
	TestConnection(ctx context.Context, ownerID, providerName, secretRef string) (bool, error)
}

// DefaultService implements Provider configuration logic.
type DefaultService struct {
	log        zerolog.Logger
	store      Store
	envService environments.Service
}

func NewService(log zerolog.Logger, store Store, envService environments.Service) *DefaultService {
	return &DefaultService{
		log:        log.With().Str("component", "ProviderConfigService").Logger(),
		store:      store,
		envService: envService,
	}
}

func (s *DefaultService) SaveConfig(ctx context.Context, cfg *ProviderConfig) error {
	return s.store.UpsertConfig(ctx, cfg)
}

func (s *DefaultService) GetConfig(ctx context.Context, id string) (*ProviderConfig, error) {
	return s.store.GetConfig(ctx, id)
}

func (s *DefaultService) ListConfigs(ctx context.Context, scope, ownerID string, params common.QueryParams) ([]*ProviderConfig, int, error) {
	return s.store.ListConfigs(ctx, scope, ownerID, params)
}

// TestConnection resolves the secretRef in any available environment and tests the connection
func (s *DefaultService) TestConnection(ctx context.Context, ownerID, providerName, secretRef string) (bool, error) {
	// For production readiness, we need to verify the secret against a real API endpoint
	// 1. Find the secret in any of the environments
	envs, _, err := s.envService.ListEnvironments(ctx, ownerID, common.QueryParams{PerPage: 1000, Page: 1})
	if err != nil {
		return false, err
	}

	var apiKey string
	for _, env := range envs {
		resolved, err := s.envService.Resolve(ctx, env.ID)
		if err == nil {
			if val, ok := resolved[secretRef]; ok && val != "" {
				apiKey = val
				break
			}
		}
	}

	if apiKey == "" {
		return false, errors.New("secret reference not found in any environment")
	}

	// 2. Perform actual HTTP ping depending on provider
	switch providerName {
	case "openai":
		req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.openai.com/v1/models", nil)
		req.Header.Set("Authorization", "Bearer "+apiKey)
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != 200 {
			return false, errors.New("failed to connect to OpenAI")
		}
	case "gemini":
		req, _ := http.NewRequestWithContext(ctx, "GET", "https://generativelanguage.googleapis.com/v1beta/models?key="+apiKey, nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != 200 {
			return false, errors.New("failed to connect to Gemini")
		}
	case "groq":
		req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.groq.com/openai/v1/models", nil)
		req.Header.Set("Authorization", "Bearer "+apiKey)
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != 200 {
			return false, errors.New("failed to connect to Groq")
		}
	case "ollama":
		// Assumes ollama is running locally for the test
		req, _ := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:11434/api/tags", nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != 200 {
			return false, errors.New("failed to connect to local Ollama")
		}
	default:
		return false, errors.New("unsupported provider")
	}

	s.log.Info().Str("provider", providerName).Msg("Provider connection tested successfully via direct HTTP")
	return true, nil
}
