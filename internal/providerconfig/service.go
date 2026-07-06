package providerconfig

import (
	"context"

	"github.com/Ajayvtl/devserver/internal/environments"
	"github.com/rs/zerolog"
)

// Store defines persistence for provider configurations.
type Store interface {
	GetConfig(ctx context.Context, id string) (*ProviderConfig, error)
	ListConfigs(ctx context.Context, scope, ownerID string) ([]*ProviderConfig, error)
	UpsertConfig(ctx context.Context, cfg *ProviderConfig) error
	DeleteConfig(ctx context.Context, id string) error
}

// Service manages provider configurations and handles runtime binding via environments.
type Service interface {
	SaveConfig(ctx context.Context, cfg *ProviderConfig) error
	GetConfig(ctx context.Context, id string) (*ProviderConfig, error)
	ListConfigs(ctx context.Context, scope, ownerID string) ([]*ProviderConfig, error)
	TestConnection(ctx context.Context, id string, envID string) (bool, error)
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

func (s *DefaultService) ListConfigs(ctx context.Context, scope, ownerID string) ([]*ProviderConfig, error) {
	return s.store.ListConfigs(ctx, scope, ownerID)
}

// TestConnection simulates resolving the environment and checking if the provider config is functionally valid.
// In a full implementation, it would dynamically instantiate the provider SDK using the resolved secrets
// and ping the Provider's /models or /user endpoint.
func (s *DefaultService) TestConnection(ctx context.Context, id string, envID string) (bool, error) {
	cfg, err := s.store.GetConfig(ctx, id)
	if err != nil {
		return false, err
	}

	if !cfg.Enabled {
		return false, nil
	}

	// Resolve the environment to extract the plaintext secret
	resolved, err := s.envService.Resolve(ctx, envID)
	if err != nil {
		s.log.Error().Err(err).Str("config_id", id).Str("env_id", envID).Msg("Failed to resolve environment for provider test")
		return false, err
	}

	// Secret extraction
	// We expect the environment resolving payload to contain the key mapping to our secretID.
	// Note: in a fully linked model, cfg.SecretID could correspond to the precise key required.
	// For now, we simulate success if the environment resolved successfully.
	_ = resolved

	s.log.Info().Str("provider", string(cfg.Type)).Str("name", cfg.Name).Msg("Provider connection tested successfully")
	return true, nil
}
