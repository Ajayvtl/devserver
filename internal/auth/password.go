package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// LocalProvider implements password-based authentication.
type LocalProvider struct{}

func (p *LocalProvider) Name() string { return "local" }

func (p *LocalProvider) Authenticate(ctx context.Context, req map[string]any) (*User, error) {
	username, _ := req["username"].(string)
	password, _ := req["password"].(string)

	if username == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	// Mock verification for WP-7.1
	if username == "admin" && password == "admin" {
		return &User{
			ID:        uuid.NewString(),
			Username:  "admin",
			Email:     "admin@devserver.local",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}, nil
	}

	return nil, ErrInvalidCredentials
}
