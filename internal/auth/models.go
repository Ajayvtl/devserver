package auth

import (
	"time"
)

// User represents an authenticated identity within the platform.
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Session tracks an active authentication session.
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Token     string    `json:"-"` // Opaque or JWT session token
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
	Revoked   bool      `json:"revoked"`
}

// TokenPair holds the active session and refresh token.
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int    `json:"expiresIn"`
}

// AuditEvent records authentication lifecycle events.
type AuditEvent struct {
	EventID   string    `json:"eventId"`
	UserID    string    `json:"userId"`
	Action    string    `json:"action"` // e.g., "login", "logout", "refresh", "revoke"
	IPAddress string    `json:"ipAddress"`
	UserAgent string    `json:"userAgent"`
	Timestamp time.Time `json:"timestamp"`
}
