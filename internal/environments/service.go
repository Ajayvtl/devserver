package environments

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
)

var (
	ErrEnvNotFound    = errors.New("environment not found")
	ErrVarNotFound    = errors.New("variable not found")
	ErrSecretNotFound = errors.New("secret not found")
)

// Store defines persistence for environments, variables, and secrets.
type Store interface {
	GetEnvironment(ctx context.Context, id string) (*Environment, error)
	ListEnvironments(ctx context.Context, ownerID string) ([]*Environment, error)
	UpsertEnvironment(ctx context.Context, env *Environment) error
	DeleteEnvironment(ctx context.Context, id string) error

	GetVariable(ctx context.Context, envID, key string) (*Variable, error)
	ListVariables(ctx context.Context, envID string) ([]*Variable, error)
	UpsertVariable(ctx context.Context, variable *Variable) error
	DeleteVariable(ctx context.Context, envID, key string) error

	GetSecret(ctx context.Context, envID, key string) (*Secret, error)
	ListSecrets(ctx context.Context, envID string) ([]*SecretReference, error)
	ListRawSecrets(ctx context.Context, envID string) ([]*Secret, error) // INTERNAL USE ONLY
	UpsertSecret(ctx context.Context, secret *Secret) error
	DeleteSecret(ctx context.Context, envID, key string) error
}

// Service manages the resolution and lifecycles of environments.
type Service interface {
	// Management
	CreateEnvironment(ctx context.Context, ownerID, name string, envType EnvType) (*Environment, error)
	ListEnvironments(ctx context.Context, ownerID string) ([]*Environment, error)

	// Variables
	SetVariable(ctx context.Context, envID, key, value string) error
	ListVariables(ctx context.Context, envID string) ([]*Variable, error)

	// Secrets
	SetSecret(ctx context.Context, envID, key, plaintext string) error
	ListSecrets(ctx context.Context, envID string) ([]*SecretReference, error)

	// Resolution
	Resolve(ctx context.Context, envID string) (map[string]string, error)
}

// DefaultService implements core environment orchestration.
type DefaultService struct {
	log    zerolog.Logger
	store  Store
	crypto CryptoService
}

func NewService(log zerolog.Logger, store Store, crypto CryptoService) *DefaultService {
	return &DefaultService{
		log:    log.With().Str("component", "EnvService").Logger(),
		store:  store,
		crypto: crypto,
	}
}

func (s *DefaultService) CreateEnvironment(ctx context.Context, ownerID, name string, envType EnvType) (*Environment, error) {
	env := &Environment{
		OwnerID: ownerID,
		Name:    name,
		Type:    envType,
	}
	if err := s.store.UpsertEnvironment(ctx, env); err != nil {
		return nil, err
	}
	return env, nil
}

func (s *DefaultService) ListEnvironments(ctx context.Context, ownerID string) ([]*Environment, error) {
	return s.store.ListEnvironments(ctx, ownerID)
}

func (s *DefaultService) SetVariable(ctx context.Context, envID, key, value string) error {
	variable := &Variable{
		EnvID: envID,
		Key:   key,
		Value: value,
	}
	return s.store.UpsertVariable(ctx, variable)
}

func (s *DefaultService) ListVariables(ctx context.Context, envID string) ([]*Variable, error) {
	return s.store.ListVariables(ctx, envID)
}

func (s *DefaultService) SetSecret(ctx context.Context, envID, key, plaintext string) error {
	ciphertext, err := s.crypto.Encrypt(plaintext)
	if err != nil {
		return err
	}

	secret := &Secret{
		EnvID: envID,
		Key:   key,
		Value: ciphertext,
	}
	return s.store.UpsertSecret(ctx, secret)
}

func (s *DefaultService) ListSecrets(ctx context.Context, envID string) ([]*SecretReference, error) {
	return s.store.ListSecrets(ctx, envID)
}

// Resolve fully expands an environment context.
// In a true fallback chain scenario, it would merge parent scopes into the child scope.
// Here we just resolve variables and decipher secrets into an ephemeral map for injection.
func (s *DefaultService) Resolve(ctx context.Context, envID string) (map[string]string, error) {
	result := make(map[string]string)

	vars, err := s.store.ListVariables(ctx, envID)
	if err != nil {
		return nil, err
	}

	for _, v := range vars {
		result[v.Key] = v.Value
	}

	rawSecrets, err := s.store.ListRawSecrets(ctx, envID)
	if err != nil {
		return nil, err
	}

	for _, sec := range rawSecrets {
		plaintext, err := s.crypto.Decrypt(sec.Value)
		if err != nil {
			s.log.Error().Err(err).Str("env_id", envID).Str("key", sec.Key).Msg("Failed to decrypt secret during environment resolution")
			return nil, err
		}
		result[sec.Key] = plaintext
	}

	return result, nil
}
