package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Ajayvtl/devserver/internal/events"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// A hardcoded secret for WP-7.1. In production, load from environment.
var jwtSecret = []byte("super-secret-key-for-devserver")

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
	ValidateToken(ctx context.Context, token string) (string, error)
	RevokeAll(ctx context.Context, userID string) error
}

// DefaultService implements the core Authentication & Identity boundaries.
type DefaultService struct {
	log       zerolog.Logger
	bus       events.Bus
	providers map[string]Provider
	store     *MySQLStore
}

func NewService(log zerolog.Logger, bus events.Bus, store *MySQLStore) *DefaultService {
	return &DefaultService{
		log:       log.With().Str("component", "AuthService").Logger(),
		bus:       bus,
		providers: make(map[string]Provider),
		store:     store,
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

	session, err := s.createSession(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	s.publishAudit(user, "login_success")

	accessToken, err := s.generateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: session.Token, // Opaque refresh token
		ExpiresIn:    3600,
	}, nil
}

func (s *DefaultService) Logout(ctx context.Context, sessionID string) error {
	session, err := s.store.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}
	if session.Revoked {
		return ErrSessionRevoked
	}

	session.Revoked = true
	if err := s.store.UpdateSession(ctx, session); err != nil {
		return err
	}

	user, _ := s.store.GetUserByID(ctx, session.UserID)
	s.publishAudit(user, "logout")

	return nil
}

func (s *DefaultService) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	session, err := s.store.GetSessionByToken(ctx, refreshToken)
	if err != nil {
		return nil, ErrSessionExpired
	}

	if session.Revoked || time.Now().UTC().After(session.ExpiresAt) {
		return nil, ErrSessionExpired
	}

	user, err := s.store.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Rotate session
	session.Revoked = true
	_ = s.store.UpdateSession(ctx, session)

	newSession, err := s.createSession(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.generateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	s.publishAudit(user, "token_refreshed")

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newSession.Token,
		ExpiresIn:    3600,
	}, nil
}

func (s *DefaultService) ValidateSession(ctx context.Context, sessionID string) (*User, error) {
	session, err := s.store.GetSession(ctx, sessionID)
	if err != nil {
		return nil, ErrSessionExpired
	}
	if session.Revoked {
		return nil, ErrSessionRevoked
	}
	if time.Now().UTC().After(session.ExpiresAt) {
		return nil, ErrSessionExpired
	}

	user, err := s.store.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *DefaultService) RevokeAll(ctx context.Context, userID string) error {
	if err := s.store.RevokeAllSessionsForUser(ctx, userID); err != nil {
		return err
	}

	user, _ := s.store.GetUserByID(ctx, userID)
	s.publishAudit(user, "revoke_all_sessions")

	return nil
}

func (s *DefaultService) ValidateToken(ctx context.Context, tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", ErrSessionExpired
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrSessionExpired
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return "", ErrSessionExpired
	}

	return userID, nil
}

func (s *DefaultService) createSession(ctx context.Context, userID string) (*Session, error) {
	sessionID := uuid.NewString()
	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		Token:     uuid.NewString(),
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().Add(1 * time.Hour).UTC(),
	}

	if err := s.store.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *DefaultService) generateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().UTC().Unix(),
		"exp": time.Now().Add(1 * time.Hour).UTC().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (s *DefaultService) publishAudit(user *User, action string) {
	var uid string
	if user != nil {
		uid = user.ID
	}

	event := AuditEvent{
		EventID:   uuid.NewString(),
		UserID:    uid,
		Action:    action,
		Timestamp: time.Now().UTC(),
	}

	// Persist to DB
	_ = s.store.SaveAuditEvent(context.Background(), event)

	// Emit audit event to central event bus
	if s.bus != nil {
		s.bus.Publish(events.AuthAuditLog, events.AuthAuditEvent{
			UserID: uid,
			Action: action,
		})
	}
}
