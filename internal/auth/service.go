package auth

import (
	"context"
	"errors"
	"time"

	"github.com/Ajayvtl/devserver/internal/events"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSessionExpired     = errors.New("session expired")
	ErrSessionRevoked     = errors.New("session revoked")
	ErrUserNotFound       = errors.New("user not found")
)

// Provider abstracts authentication logic (e.g., Local Password, OAuth).
type Provider interface {
	Name() string
	Authenticate(ctx context.Context, req map[string]any) (*User, error)
}

// Service manages the authentication lifecycle.
type Service interface {
	Login(ctx context.Context, providerName string, credentials map[string]any) (*TokenPair, error)
	Logout(ctx context.Context, sessionID string) error
	Refresh(ctx context.Context, refreshToken string) (*TokenPair, error)
	ValidateSession(ctx context.Context, sessionID string) (*User, error)
	RevokeAll(ctx context.Context, userID string) error
}

// DefaultService implements the core Authentication & Identity boundaries.
type DefaultService struct {
	log       zerolog.Logger
	bus       events.Bus
	providers map[string]Provider
	// In-memory stores for WP-7.1. Later migrated to DB.
	users    map[string]*User
	sessions map[string]*Session
}

func NewService(log zerolog.Logger, bus events.Bus) *DefaultService {
	return &DefaultService{
		log:       log.With().Str("component", "AuthService").Logger(),
		bus:       bus,
		providers: make(map[string]Provider),
		users:     make(map[string]*User),
		sessions:  make(map[string]*Session),
	}
}

func (s *DefaultService) RegisterProvider(p Provider) {
	s.providers[p.Name()] = p
}

func (s *DefaultService) Login(ctx context.Context, providerName string, credentials map[string]any) (*TokenPair, error) {
	provider, exists := s.providers[providerName]
	if !exists {
		return nil, errors.New("authentication provider not found")
	}

	user, err := provider.Authenticate(ctx, credentials)
	if err != nil {
		s.publishAudit(user, "login_failed")
		return nil, err
	}

	// Upsert user (In-memory mock for WP-7.1)
	s.users[user.ID] = user

	session := s.createSession(user.ID)

	s.publishAudit(user, "login_success")

	return &TokenPair{
		AccessToken:  session.Token,
		RefreshToken: uuid.NewString(), // Mock refresh token
		ExpiresIn:    3600,
	}, nil
}

func (s *DefaultService) Logout(ctx context.Context, sessionID string) error {
	session, exists := s.sessions[sessionID]
	if !exists || session.Revoked {
		return ErrSessionRevoked
	}
	session.Revoked = true

	user := s.users[session.UserID]
	s.publishAudit(user, "logout")

	return nil
}

func (s *DefaultService) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	// Simple mock implementation for WP-7.1
	return nil, errors.New("refresh flow requires persistent store (upcoming)")
}

func (s *DefaultService) ValidateSession(ctx context.Context, sessionID string) (*User, error) {
	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, ErrSessionExpired
	}
	if session.Revoked {
		return nil, ErrSessionRevoked
	}

	user, exists := s.users[session.UserID]
	if !exists {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *DefaultService) RevokeAll(ctx context.Context, userID string) error {
	for _, session := range s.sessions {
		if session.UserID == userID {
			session.Revoked = true
		}
	}

	user := s.users[userID]
	s.publishAudit(user, "revoke_all_sessions")

	return nil
}

func (s *DefaultService) createSession(userID string) *Session {
	sessionID := uuid.NewString()
	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		Token:     uuid.NewString(),
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().Add(1 * time.Hour).UTC(),
	}
	s.sessions[sessionID] = session
	return session
}

func (s *DefaultService) publishAudit(user *User, action string) {
	var uid string
	if user != nil {
		uid = user.ID
	}

	// Emit audit event to central event bus
	if s.bus != nil {
		s.bus.Publish(events.AuthAuditLog, events.AuthAuditEvent{
			UserID: uid,
			Action: action,
		})
	}
}
