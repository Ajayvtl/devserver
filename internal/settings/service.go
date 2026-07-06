package settings

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
)

var (
	ErrSettingNotFound     = errors.New("setting not found")
	ErrIntegrationNotFound = errors.New("integration not found")
)

// Store defines persistent storage for settings and integrations.
type Store interface {
	GetSetting(ctx context.Context, scope Scope, ownerID, key string) (*Setting, error)
	ListSettings(ctx context.Context, scope Scope, ownerID string) ([]*Setting, error)
	UpsertSetting(ctx context.Context, setting *Setting) error
	DeleteSetting(ctx context.Context, scope Scope, ownerID, key string) error

	GetIntegration(ctx context.Context, id string) (*Integration, error)
	ListIntegrations(ctx context.Context, scope Scope, ownerID string) ([]*Integration, error)
	UpsertIntegration(ctx context.Context, integration *Integration) error
	DeleteIntegration(ctx context.Context, id string) error
}

// Service manages the business logic for settings and integrations.
type Service interface {
	GetSettingValue(ctx context.Context, scope Scope, ownerID, key string) (string, error)
	ListSettings(ctx context.Context, scope Scope, ownerID string) ([]*Setting, error)
	SaveSetting(ctx context.Context, scope Scope, ownerID, key, value string) error

	GetIntegration(ctx context.Context, id string) (*Integration, error)
	ListIntegrations(ctx context.Context, scope Scope, ownerID string) ([]*Integration, error)
	SaveIntegration(ctx context.Context, integration *Integration) error
}

// DefaultService implements core settings logic.
type DefaultService struct {
	log   zerolog.Logger
	store Store
}

func NewService(log zerolog.Logger, store Store) *DefaultService {
	return &DefaultService{
		log:   log.With().Str("component", "SettingsService").Logger(),
		store: store,
	}
}

func (s *DefaultService) GetSettingValue(ctx context.Context, scope Scope, ownerID, key string) (string, error) {
	setting, err := s.store.GetSetting(ctx, scope, ownerID, key)
	if err != nil {
		return "", err
	}
	return setting.Value, nil
}

func (s *DefaultService) ListSettings(ctx context.Context, scope Scope, ownerID string) ([]*Setting, error) {
	return s.store.ListSettings(ctx, scope, ownerID)
}

func (s *DefaultService) SaveSetting(ctx context.Context, scope Scope, ownerID, key, value string) error {
	setting := &Setting{
		Scope:   scope,
		OwnerID: ownerID,
		Key:     key,
		Value:   value,
	}
	return s.store.UpsertSetting(ctx, setting)
}

func (s *DefaultService) GetIntegration(ctx context.Context, id string) (*Integration, error) {
	return s.store.GetIntegration(ctx, id)
}

func (s *DefaultService) ListIntegrations(ctx context.Context, scope Scope, ownerID string) ([]*Integration, error) {
	return s.store.ListIntegrations(ctx, scope, ownerID)
}

func (s *DefaultService) SaveIntegration(ctx context.Context, integration *Integration) error {
	return s.store.UpsertIntegration(ctx, integration)
}
